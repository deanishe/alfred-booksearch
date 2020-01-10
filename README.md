
<div align="center">
    <img height="128" width="128" src="https://git.deanishe.net/deanishe/alfred-goodreads/raw/branch/master/icons/icon.png">
</div>


Goodreads Book Search for Alfred 4
==================================

Search for movies, TV shows and people in [Alfred 4][alfred].

- [Download & installation](#download--installation)
- [Usage](#usage)
- [Configuration](#configuration)
  - [Custom URLs](#custom-urls)
- [Licensing & thanks](#licensing--thanks)

Download & installation
-----------------------

[Grab the workflow from the releases page][download]. Download the
`Goodreads-Search-X.X.X.alfredworkflow` file and double-click it to install.

Usage
-----

- `gr <query>` — Search for a book.
    - `↩` — Open book details in browser.
    - `⌘↩` — Open author page in browser.
    - `...` — Open custom URL (see [configuration](#configuration)).


Configuration
-------------

There are a couple of options in the workflow's [configuration sheet][confsheet] (the `[x]` button in Alfred Preferences):

| Setting         | Description                          |
| --------------- | ------------------------------------ |
| `API_KEY`       | Goodreads.com API key                |
| `MAX_CACHE_AGE` | How long to cache search results for |


### Custom URLs

You can assign custom URLs to arbitrary modifier keys by setting
`URL_MOD[_MOD...]` variables. So `URL_OPT` sets a URL to open
with you press `⌥↩` on a result, and `URL_OPT_SHIFT` sets a
URL to open when you press `⌥⇧↩` on a result.

You can use shell-like variable expansion in URLs to insert
book or author information in to the URL, e.g.
`https://duckduckgo.com/?q=${Title}`.

The following variables are available:

| Variable      | Meaning                      |
| ------------- | ---------------------------- |
| `Title`       | Book title                   |
| `Author`      | Author name                  |
| `AuthorID`    | Author ID (for Goodreads)    |
| `AuthorURL`   | Author page on Goodreads.com |
| `AuthorTitle` | Author name & book title     |
| `Year`        | Year of first publication    |
| `Rating`      | Average rating (0.0 – 5.0)   |
| `URL`         | Book page on Goodreads.com   |
| `ImageURL`    | URL of cover image           |

Values are URL query-escaped by default (i.e. spaces are replaced with
`+`). If you need path-escaped values (i.e. spaces are replaced with
`%20`), use the `*Alt` variants (`TitleAlt`, `AuthorAlt` etc.). The
original, unescaped values are available via the `*Raw` variables
(`TitleRaw`, `AuthorRaw` etc.).


Licensing & thanks
------------------

This workflow is released under the [MIT Licence][mit].

It is based on [AwGo][awgo] ([MIT][mit]).


The icons are from or based on [Font Awesome][awesome] ([SIL][sil]).


[alfred]: https://alfredapp.com/
[confsheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
[awgo]: https://github.com/deanishe/awgo
[download]: https://git.deanishe.net/deanishe/alfred-imdb/releases/latest
[issues]: https://git.deanishe.net/deanishe/alfred-imdb/issues
[sil]: http://scripts.sil.org/cms/scripts/page.php?site_id=nrsi&id=OFL
[mit]: https://opensource.org/licenses/MIT
[awesome]: http://fortawesome.github.io/Font-Awesome/
