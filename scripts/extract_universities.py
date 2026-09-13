"""Extract unique (name, country) pairs from NUS SEP placement PDFs."""

import argparse
import csv
from pathlib import Path

import pdfplumber


def clean(value):
    return " ".join((value or "").split())


def extract(pdf_path):
    records = []
    with pdfplumber.open(pdf_path) as pdf:
        for number, page in enumerate(pdf.pages, 1):
            found = False
            for table in page.extract_tables():
                columns = None
                for row in table:
                    cells = [clean(cell) for cell in row]
                    if "Partner University" in cells and "Region" in cells:
                        columns = (cells.index("Partner University"), cells.index("Region"))
                        found = True
                        continue
                    if columns is None or not any(cells):
                        continue
                    if len(cells) <= max(columns):
                        raise ValueError(f"{pdf_path.name}, page {number}: incomplete row")
                    name, country = (cells[index] for index in columns)
                    if not name or not country:
                        raise ValueError(f"{pdf_path.name}, page {number}: missing name/country")
                    records.append((name, country))
            if not found:
                raise ValueError(f"{pdf_path.name}, page {number}: table header not found")
    if not records:
        raise ValueError(f"No universities found in {pdf_path.name}")
    return records


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("pdfs", nargs="+", type=Path)
    parser.add_argument("--output", required=True, type=Path)
    args = parser.parse_args()
    records = set()
    for path in args.pdfs:
        extracted = extract(path)
        print(f"{path.name}: extracted {len(extracted)} rows")
        records.update(extracted)
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("w", newline="", encoding="utf-8") as output:
        writer = csv.writer(output)
        writer.writerow(["name", "country"])
        writer.writerows(sorted(records, key=lambda row: (row[1], row[0])))
    print(f"Wrote {len(records)} unique pairs to {args.output}")


if __name__ == "__main__":
    main()
