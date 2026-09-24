"""Remove listed duplicates and add local universities to university-refined.csv."""

import csv
from pathlib import Path


# Update this list when more duplicate entries are identified in universities.csv.
TO_REMOVE = {
    ("Queen's University, Kingston", "Canada"),
    ("University Lausanne", "Switzerland"),
    ("University of Illinois at Urbana Champaign", "USA"),
    ("University of Wisconsin - Madison", "USA"),
    ("The University of Bath", "United Kingdom"),
    ("The University of Nottingham", "United Kingdom"),
}

# Insert these names exactly as written, including the acronyms.
MANUAL_UNIVERSITIES = {
    ("Nanyang Technological University (NTU)", "Singapore"),
    ("Singapore Management University (SMU)", "Singapore"),
    ("Singapore University of Social Sciences (SUSS)", "Singapore"),
    ("Singapore University of Technology & Design (SUTD)", "Singapore"),
    ("Singapore Institute of Technology (SIT)", "Singapore"),
}


def main():
    folder = Path(__file__).resolve().parents[1] / "data" / "generated"
    with (folder / "universities.csv").open(newline="", encoding="utf-8-sig") as source:
        reader = csv.reader(source)
        header = next(reader)
        universities = {tuple(row) for row in reader}

    refined = (universities - TO_REMOVE) | MANUAL_UNIVERSITIES
    with (folder / "university-refined.csv").open("w", newline="", encoding="utf-8") as output:
        writer = csv.writer(output)
        writer.writerow(header)
        writer.writerows(sorted(refined, key=lambda row: (row[1], row[0])))


if __name__ == "__main__":
    main()
