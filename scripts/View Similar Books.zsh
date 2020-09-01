#!/bin/zsh -e

eval "$( ./alfred-booksearch -export )"

if [[ "$WORK_ID" == "0" ]]; then
	echo "WORK_ID is not set" >&2
	exit 1
fi

/usr/bin/open "https://www.goodreads.com/book/similar/${WORK_ID}"
