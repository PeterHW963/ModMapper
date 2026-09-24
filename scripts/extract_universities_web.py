"""Extract name,country from a saved NUS GRO partner directory HTML file."""

import argparse
import csv
from pathlib import Path
import re

from bs4 import BeautifulSoup


def clean(value):
    return " ".join(value.split())


def extract(html):
    soup = BeautifulSoup(html, "html.parser")
    records = set() # deduplicate some entries
    for heading in soup.select("h4.inner-header"):
        country = clean(heading.get_text(" ", strip=True))
        section = heading.find_next_sibling("div", class_="show-more")
        for item in section.select("li"):
            for annotation in item.select("sup"):
                annotation.decompose()
            name = clean(item.get_text(" ", strip=True))
            # remove any notations on the uni name
            name = re.sub(r"\s*\(formerly\b[^)]*\)", "", name, flags=re.I)
            name = re.sub(r"\s*[-–—]\s*on hold\s*$", "", name, flags=re.I)
            if name:
                records.add((name, country))
    return sorted(records, key=lambda row: (row[1], row[0]))


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--html", type=Path, default=Path("data/raw/nus-partner-universities.html"))
    parser.add_argument("--output", type=Path, default=Path("data/generated/universities.csv"))
    args = parser.parse_args()
    records = extract(args.html.read_bytes())
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("w", newline="", encoding="utf-8") as output:
        writer = csv.writer(output)
        writer.writerow(["name", "country"])
        writer.writerows(records)
    print(f"Wrote {len(records)} unique name/country pairs to {args.output}")


if __name__ == "__main__":
    main()
