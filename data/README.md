# Instructions to save raw data:

## 1. Save the source page manually

Open the [NUS partner directory](https://www.nus.edu.sg/gro/global-programmes/student-exchange/partner-universities)
in your browser.

Use **Save page as** (usually Ctrl+S), choose an HTML format (not PDF or MHTML),
and save the page to:

```text
\data\raw\nus-partner-universities.html
```

## 2. Install the HTML parsing dependency

`scripts/extract_universities_web.py` converts the HTML into CSV and
`scripts/requirements-web.txt` documents Beautiful Soup dependency requirement.

From project root, create the virtual environment if it does not already exist:

```powershell
python -m venv .venv
```

Install dependencies into it:

```powershell
.\.venv\Scripts\python.exe -m pip install -r scripts/requirements-web.txt
```

## 3. Extract the saved HTML

From project root, run the extractor on the saved HTML (the default input is the path below; no network request is made):

```powershell
python .\scripts\extract_universities_web.py --html .\data\raw\nus-partner-universities.html --output .\data\generated\universities.csv
```

Input is by default `data/raw/nus-partner-universities.html` and Output is by default `data/generated/universities.csv` with exactly `name,country` columns.

If you choose to name the files in some other way, you may do the following:

```powershell
python .\scripts\extract_universities_web.py --html <input filepath> --output <output filepath>
```

## 4. Import the reviewed CSV

The importer defaults to a different location, so pass the path explicitly from
`<root>\server`:

```powershell
go run ./cmd/import-universities -csv ../data/generated/universities.csv -env ../.env
```

or if you placed the genearated csv file and env files differently:

```powershell
go run ./cmd/import-universities -csv <csv filepath> -env <env filepath>
```

