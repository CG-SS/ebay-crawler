# Ebay crawler

Crawler written in Go that craws Ebay and extract item information such as title, price, condition and the URL using [GoQuery](https://github.com/PuerkitoBio/goquery). It also uses [argparse](https://github.com/akamensky/argparse) for parsing command-line arguments.

Features:

- Asynchronously writes the JSON files to disk.
- Multiple crawler workers
- Supports pagination
- Supports crawling by item condition

### Structure



 ```bash
 ├───.idea
├───cmd
│   └───local
└───internal
    ├───config
    ├───crawler
    ├───itemcondition
    ├───model
    └───persistence
 ```

- `cmd/local`: contains the application entrypoint
- `internal/config`: contains the app config
- `internal/crawler`: contains the main business logic
- `internal/itemcondition`: contains the item condition enum
- `internal/mode`: contains the data model for the app
- `internal/persistence`: contains the business logic for persistence
