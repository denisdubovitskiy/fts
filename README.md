# FTS

This tool is made for quick searching using SQLite FTS5.

   `WARNING: This tool is made "just for fun" without any warranty. Use at your own risk.`

1. Build (you need the go compiler)

    ```shell
    make build
    ```

2. Install

    ```shell
    sudo install ./bin/fts /usr/local/bin
    ```

3. Index your documents

   ```shell
   fts index /path/to/documents
   ```

4. Query via CLI (you can use any suitable query for Sqlite FTS5)

   ```shell
   fts query "weather AND (sunny OR cloudy)"
   ```

5. Serve an embedded web interface

   ```shell
   fts web /path/to/documents
   ```