
Configuration
=============

The workflow offers a few options and links under the `bkconf` keyword, but it's primarily configured via its [configuration sheet][confsheet] (the `[x]` button in Alfred Preferences).


<!-- MarkdownTOC autolink="true" bracket="round" levels="1,2,3" autoanchor="true" -->

- [Inline configuration](#inline-configuration)
- [Workflow configuration sheet](#workflow-configuration-sheet)
- [Adding custom actions](#adding-custom-actions)

<!-- /MarkdownTOC -->


<a id="inline-configuration"></a>
Inline configuration
--------------------

Keyword: `bkconf`

<a id="workflow-is-up-to-dateworkflow-update-available"></a>
#### Workflow Update Available / Workflow Is Up To Date ###

If a newer version of the workflow is available, "Workflow Update Available" is shown. Action this item to update the workflow.


#### Workflow Authorised / Workflow Not Authorised ####

Whether you've authorised the workflow to access your Goodreads account via OAuth. If "Workflow Not Authorised" is shown, action the item to go to goodreads.com and authorise the workflow.

You can deauthorise the workflow (i.e. delete the OAuth tokens) with `⌘↩`.


#### Open Scripts Folder ####

Open custom scripts folder (see [scripts][scripts] for details).


#### Open Docs ####

Open this documentation in your browser.


#### Get Help ####

Open workflow issue tracker in your browser.


#### Report Bug ####

Open workflow issue tracker in your browser.


<a id="workflow-configuration-sheet"></a>
Workflow configuration sheet
----------------------------

There are only a couple of configuration options by default, but you can add more to customise the workflow.

|     Variable     |  Default Value   |                                               Description                                               |
|------------------|------------------|---------------------------------------------------------------------------------------------------------|
| `ACTION_DEFAULT` | `Open Book Page` | The script run when you press `↩` on a book item.                                                       |
| `ACTION_ALT`     | `View Series`    | The script run when you press `⌥↩` on a book item.                                                      |
| `EXPORT_DETAILS` | `false`          | Whether all book details should be fetched before running a script (see [scripts][scripts] for details) |
| `USER_ID`        |                  | Your Goodreads ID. Saved by the workflow when you log in.                                               |
| `USER_NAME`      |                  | Your Goodreads username. Saved by the workflow when you log in.                                         |


<a id="adding-custom-actions"></a>
Adding custom actions
---------------------

The actions you can perform on books are primarily defined by scripts. Several scripts are included with the workflow (in the `scripts` subdirectory of the workflow's folder), and you can also add your own (see [scripts][scripts]).

You can assign any built-in or custom script a hotkey via the configuration by creating a variable with a name of the form `ACTION_<KEY>` and the name of the script (without file extension) as its value. For example:

|        Variable       |        Value         |                                     Description                                      |
|-----------------------|----------------------|--------------------------------------------------------------------------------------|
| `ACTION_CTRL`         | `Mark as Read`       | Run the built-in `Mark as Read.zsh` script when you press `^↩` on a book item        |
| `ACTION_CTRL_SHIFT`   | `View Similar Books` | Run the built-in `View Similar Books.zsh` script when you press `^⇧↩` on a book item |
| `ACTION_CMD_OPT_CTRL` | `My Custom Script`   | Run your custom script named `My Custom Script` when you press `^⌥⌘↩` on a book item |

You can combine modifier keys arbitrarily.

To simplify adding custom hotkeys, hitting `⌘C` on an action in the "All Actions…" list (`⌘↩` on a book item) will copy its name to the clipboard for easy pasting into the configuration sheet.

[↑ Documentation][top]

[top]: ./README.md
[scripts]: ./scripts.md
[confsheet]: https://www.alfredapp.com/help/workflows/advanced/variables/#environment
