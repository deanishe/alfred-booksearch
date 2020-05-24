#!/bin/zsh -e

echo -n "https://www.goodreads.com/book/show/${BOOK_ID}" | /usr/bin/pbcopy

./alfred-booksearch -beep
