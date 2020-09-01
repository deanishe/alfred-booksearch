Scripts
=======

The actions you can perform on book results in the workflow are defined by scripts. The workflow includes several built-in scripts for common actions, but you can also define your own.

You can also assign hotkeys to scripts, so you can run the script directly from a book item (see [configuration][configuration]).

<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3,4" autoanchor="true" -->

- [Built-in scripts/actions](#built-in-scriptsactions)
- [Writing custom scripts](#writing-custom-scripts)
    - [Environment variables](#environment-variables)
        - [Default variables](#default-variables)
        - [Details variables](#details-variables)
        - [Formatted versions](#formatted-versions)
    - [JSON](#json)
- [Helper functions](#helper-functions)
- [Script icons](#script-icons)
- [Example script](#example-script)

<!-- /MarkdownTOC -->

<a id="built-in-scriptsactions"></a>
Built-in scripts/actions
------------------------

The workflow includes the following scripts (in the internal `scripts` subdirectory):

|           Script           |                  Description                   |
|----------------------------|------------------------------------------------|
| `Add to Currently Reading` | Add book to your "Currently Reading" bookshelf |
| `Add to Shelves`           | Add book to one or more shelves                |
| `Add to Want to Read`      | Add book to your "Want to Read" bookshelf      |
| `Copy Goodreads Link`      | Copy URL of book's page on goodreads.com       |
| `Mark as Read`             | Add book to your "Read" bookshelf              |
| `Open Author Page`         | Open author's page on goodreads.com            |
| `Open Book Page`           | Open book's page on goodreads.com              |
| `View Author’s Books`      | View list of author's books in Alfred          |
| `View Series`              | View all books in a book's series in Alfred    |
| `View Similar Books`       | Open list of similar books on goodreads.com    |


<a id="writing-custom-scripts"></a>
Writing custom scripts
----------------------

You can add new actions to the workflow by saving custom scripts in the user scripts directory. Use `bkconf` > `Open Scripts Folder` to open this folder in Finder. **Do not put your own scripts in the workflow's internal `scripts` directory: they'll be overwritten/deleted when the workflow is updated.**

A "script" may be an executable file or any of [the script types understood by AwGo][script-types]. Its base name (i.e. without file extension) will be the name used by the workflow. If a custom script has the same name as a built-in one, it will override the built-in.


<a id="environment-variables"></a>
### Environment variables ###

When the workflow executes a script, it passes book properties via environment variables. Unfortunately, different parts of the Goodreads API provide different levels of detail about books, so while some variables are always available, others require the workflow to fetch the book's details from the API. As this can take a few seconds if the details aren't already cached, it isn't done by default (though you can force that behaviour by setting `EXPORT_DETAILS` to `true` in the workflow's configuration sheet). So if your script requires more than basic info about the book, it must call the workflow binary to retrieve it.

This is done by calling `./alfred-booksearch -export` (all scripts are run with the workflow's directory as the working directory).

In a shell script, you can export book properties to environment variables with:

```bash
eval "$( ./alfred-booksearch -export )"
```

In other languages, you can get book properties as JSON with:

```bash
./alfred-booksearch -export -json
```

<a id="default-variables"></a>
#### Default variables ####

These variables are always available (provided the book has corresponding properties):

|      Variable     |                   Description                    |
|-------------------|--------------------------------------------------|
| `BOOK_ID`         | Goodreads ID of the book                         |
| `BOOK_URL`        | URL of book's page on goodreads.com              |
| `TITLE`           | The book title                                   |
| `TITLE_NO_SERIES` | Book title without series info                   |
| `SERIES`          | Title of series book is part of (if it's in one) |
| `AUTHOR`          | Name of the author                               |
| `AUTHOR_ID`       | Author's Goodreads ID                            |
| `AUTHOR_URL`      | URL of author's page on goodreads.com            |
| `YEAR`            | Year book was published (often not available)    |
| `RATING`          | Book rating (0.0–5.0)                            |
| `IMAGE_URL`       | URL of book's cover (often not available)        |
| `USER_ID`         | Your Goodreads user ID                           |
| `USER_NAME`       | Your Goodreads username                          |


<a id="details-variables"></a>
#### Details variables ####

The following variables are only available if you export the books details:

|        Variable        |            Description            |
|------------------------|-----------------------------------|
| `DESCRIPTION`          | Plaintext description of the book |
| `DESCRIPTION_HTML`     | HTML description of the book      |
| `DESCRIPTION_MARKDOWN` | Markdown description of book      |
| `SERIES_ID`            | Series' Goodreads ID              |
| `ISBN`                 | Book's ISBN                       |
| `ISBN13`               | Book's ISBN 13                    |


<a id="formatted-versions"></a>
#### Formatted versions ####

Two additional, URL-escaped variants of each of the above variables are also exported to make it easier to insert them into URLs.

Add the suffix `_QUOTED` for a path-escaped (i.e. spaces are replaced with `%20`) version of the variable, or the suffix `_QUOTED_PLUS` for a query-escaped version (i.e. spaces are replaced with `+`).

For example, `TITLE_QUOTED` is the path-escaped book's title, and `TITLE_QUOTED_PLUS` is the query-escaped book's title.


<a id="json"></a>
### JSON ###

The JSON emitted by `./alfred-booksearch -export -json` has the following format:

```json
{
  "ID": 23106013,
  "WorkID": 42654036,
  "ISBN": "0593199308",
  "ISBN13": "9780593199305",
  "Title": "Battle Ground (The Dresden Files, #17)",
  "TitleNoSeries": "Battle Ground",
  "Series": {
    "Title": "The Dresden Files",
    "Position": 17,
    "ID": 40346,
    "Books": null
  },
  "Author": {
    "ID": 10746,
    "Name": "Jim Butcher",
    "URL": "https://www.goodreads.com/author/show/10746"
  },
  "PubDate": "2020-09-29",
  "Rating": 4.46,
  "Description": "THINGS ARE ABOUT TO GET SERIOUS FOR HARRY DRESDEN, CHICAGO’S ONLY PROFESSIONAL WIZARD, in the next entry in the #1 New York Times bestselling Dresden Files.<br /><br />Harry has faced terrible odds before. He has a long history of fighting enemies above his weight class. The Red Court of vampires. The fallen angels of the Order of the Blackened Denarius. The Outsiders.<br /><br />But this time it’s different. A being more powerful and dangerous on an order of magnitude beyond what the world has seen in a millennium is coming. And she’s bringing an army. The Last Titan has declared war on the city of Chicago, and has come to subjugate humanity, obliterating any who stand in her way.<br /><br />Harry’s mission is simple but impossible: Save the city by killing a Titan. And the attempt will change Harry’s life, Chicago, and the mortal world forever.",
  "URL": "https://www.goodreads.com/book/show/23106013",
  "ImageURL": "https://i.gr-assets.com/images/S/compressed.photo.goodreads.com/books/1587778549l/23106013._SX98_.jpg"
}

```

<a id="helper-functions"></a>
Helper functions
----------------

In addition to `-export`, the workflow binary also provides some additional helper functions to allow you to perform workflow actions or direct its behaviour. For example, you can run `./alfred-booksearch -hide=(true|false)` to tell the workflow whether to hide Alfred's window after your script is executed.

**Note: You may only call `alfred-booksearch` _once_ in your script with any of the following options because it communicates with the workflow via JSON.**

|                  Flag                  |                                                 Description                                                  |
|----------------------------------------|--------------------------------------------------------------------------------------------------------------|
| `-action <name>`                       | Set next action for workflow to run, e.g. `-action search` to open the book search after running your script |
| `-add <shelf>...`                      | Add book to named shelves, e.g. `-add to-read`                                                               |
| `-beep`                                | Play "morse" sound effect                                                                                    |
| `-notify <title> [-message <message>]` | Show a notification                                                                                          |
| `-passvars=true/false`                 | Pass workflow variables to next action                                                                       |
| `-hide=true/false`                     | Hide/show Alfred after script is run                                                                         |
| `-query <query>`                       | Set query to pass to next action                                                                             |



<a id="script-icons"></a>
Script icons
------------

You can assign an icon to a script by putting an icon with the same base name (i.e. without file extension) as the script in the user scripts directory. Supported icon types are PNG, JPG, GIF, ICNS.


<a id="example-script"></a>
Example script
--------------

As an example, here's how to add a "Search on Amazon.com" script. (Check out the built-in script for more examples.)

1. Enter `bkconf scripts` into Alfred, and action the "Open Scripts Directory" item.
2. Create a new file called "Search on Amazon.com.zsh" in the folder.
3. Add the following to the script:

```bash
# use _QUOTED_PLUS variants, as that's the format Amazon's search URL uses
url="https://www.amazon.com/s?i=stripbooks&k=${TITLE_NO_SERIES_QUOTED_PLUS}+${AUTHOR_QUOTED_PLUS}"
# open the URL in default browser
/usr/bin/open "${url}"
# ensure Alfred's window closes
./alfred-booksearch -hide
```

There's no need to make the script executable, as the workflow knows to run .zsh files with /bin/zsh.

4. Give your script a custom icon by saving an Amazon icon (such as [this one][amazon-icon]) to the same directory with the name `Search on Amazon.com.png`.
5. Search for a book (keyword `bk`), then use `⌘↩` to show all actions. You should see `Search on Amazon.com` in there.
6. Select the `Search on Amazon.com` action and hit `⌘C` to copy its name to the clipboard.
7. Open the workflow in Alfred Preferences and then open its configuration sheet (the `[x]` icon).
8. Add a new variable called `ACTION_SHIFT`, place the cursor in the Value cell and press `⌘V` to paste the clipboard contents.
9. Click "Save" (or press `⌘S`)

You should now be able to search for a book on amazon.com by hitting `⇧↩` on a book in the workflow's search results.

[↑ Documentation][top]

[top]: ./README.md
[configuration]: ./configuration.md
[script-types]: https://godoc.org/github.com/deanishe/awgo/util#Runner
[amazon-icon]: https://github.com/deanishe/alfred-searchio/raw/000243deca20c79024d27a50c7b301c44a5de4a9/src/icons/engines/amazon.png
