// Copyright (c) 2020 Dean Jackson <deanishe@deanishe.net>
// MIT Licence applies http://opensource.org/licenses/MIT
// Created on 2020-07-18

package gr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTML2Markdown(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, in, x string
	}{
		{"empty string", "", ""},
		{"four dots", "....", "…"},
		{"four dots with spaces", ". . . .", "…"},
		{"three dots with spaces", ". . .", "…"},
		{"three dots", "...", "…"},
		{"dresden", `<b>HARRY DRESDEN — WIZARD</b><br /><br /><i>Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.</i><br /><br />Harry Dresden is the best at what he does. Well, technically, he's the <i>only</i> at what he does. So when the Chicago P.D. has a case that transcends mortal creativity or capability, they come to him for answers. For the "everyday" world is actually full of strange and magical things—and most don't play well with humans. That's where Harry comes in. Takes a wizard to catch a—well, whatever. There's just one problem. Business, to put it mildly, stinks.<br /><br />So when the police bring him in to consult on a grisly double murder committed with black magic, Harry's seeing dollar signs. But where there's black magic, there's a black mage behind it. And now that mage knows Harry's name. And that's when things start to get interesting.<br /><br />Magic - it can get a guy killed.`,
			`**HARRY DRESDEN — WIZARD**

*Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.*

Harry Dresden is the best at what he does. Well, technically, he's the *only* at what he does. So when the Chicago P.D. has a case that transcends mortal creativity or capability, they come to him for answers. For the "everyday" world is actually full of strange and magical things—and most don't play well with humans. That's where Harry comes in. Takes a wizard to catch a—well, whatever. There's just one problem. Business, to put it mildly, stinks.

So when the police bring him in to consult on a grisly double murder committed with black magic, Harry's seeing dollar signs. But where there's black magic, there's a black mage behind it. And now that mage knows Harry's name. And that's when things start to get interesting.

Magic - it can get a guy killed.`},
		{"dresden with link", `<i>An alternative cover edition with a different page count exists <a href="https://www.goodreads.com/book/show/13511897.here" title="here" rel="nofollow">here</a>.</i><br /><br />Harry Dresden - Wizard<br />Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.<br /><br />Harry Dresden has faced some pretty terrifying foes during his career. Giant scorpions. Oversexed vampires. Psychotic werewolves. It comes with the territory when you're the only professional wizard in the Chicago-area phone book.<br /><br />But in all Harry's years of supernatural sleuthing, he's never faced anything like this: The spirit world has gone postal. All over Chicago, ghosts are causing trouble - and not just of the door-slamming, boo-shouting variety. These ghosts are tormented, violent, and deadly. Someone - or <i>something</i> - is purposely stirring them up to wreak unearthly havoc. But why? And why do so many of the victims have ties to Harry? If Harry doesn't figure it out soon, he could wind up a ghost himself....`,
			`*An alternative cover edition with a different page count exists here.*

Harry Dresden - Wizard  
Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.

Harry Dresden has faced some pretty terrifying foes during his career. Giant scorpions. Oversexed vampires. Psychotic werewolves. It comes with the territory when you're the only professional wizard in the Chicago-area phone book.

But in all Harry's years of supernatural sleuthing, he's never faced anything like this: The spirit world has gone postal. All over Chicago, ghosts are causing trouble - and not just of the door-slamming, boo-shouting variety. These ghosts are tormented, violent, and deadly. Someone - or *something* - is purposely stirring them up to wreak unearthly havoc. But why? And why do so many of the victims have ties to Harry? If Harry doesn't figure it out soon, he could wind up a ghost himself…`},
		{"dresden many links", `Here, together for the first time, are the shorter from Jim Butcher's DRESDEN FILES series — a compendium of cases that Harry and his cadre of allies managed to close in record time. The tales range from the deadly serious to the absurdly hilarious. Also included is a new, never-before-published novella that takes place after the cliff-hanger ending of the new April 2010 hardcover, <i>Changes</i>.<br /><br />Contains:<br />+ "Restoration of Faith"<br />+ "Vignette"<br />+ "Something Borrowed" -- from <em>
  <a href="https://www.goodreads.com/book/show/84156.My_Big_Fat_Supernatural_Wedding" title="My Big Fat Supernatural Wedding" rel="nofollow">My Big Fat Supernatural Wedding</a>
</em><br />+ "It's My Birthday Too" -- from <em>
  <a href="https://www.goodreads.com/book/show/140098.Many_Bloody_Returns" title="Many Bloody Returns" rel="nofollow">Many Bloody Returns</a>
</em><br />+ "Heorot" -- from <em>
  <a href="https://www.goodreads.com/book/show/1773616.My_Big_Fat_Supernatural_Honeymoon" title="My Big Fat Supernatural Honeymoon" rel="nofollow">My Big Fat Supernatural Honeymoon</a>
</em><br />+ "Day Off" -- from <em>
  <a href="https://www.goodreads.com/book/show/2871256.Blood_Lite" title="Blood Lite" rel="nofollow">Blood Lite</a>
</em><br />+ "<a href="https://www.goodreads.com/book/show/2575572.Backup" title="Backup" rel="nofollow">Backup</a>" -- novelette from Thomas' point of view, originally published by Subterranean Press<br />+ "The Warrior" -- novelette from <em>
  <a href="https://www.goodreads.com/book/show/3475145.Mean_Streets" title="Mean Streets" rel="nofollow">Mean Streets</a>
</em><br />+ "Last Call" -- from <em>
  <a href="https://www.goodreads.com/book/show/6122181.Strange_Brew" title="Strange Brew" rel="nofollow">Strange Brew</a>
</em><br />+ "Love Hurts" -- from <em>
  <a href="https://www.goodreads.com/book/show/7841656.Songs_of_Love_and_Death" title="Songs of Love and Death" rel="nofollow">Songs of Love and Death</a>
</em><br />+ <em>Aftermath</em> -- all-new novella from Murphy's point of view, set forty-five minutes after the end of <em>
  <a href="https://www.goodreads.com/book/show/6585201.Changes" title="Changes" rel="nofollow">Changes</a>
</em>`,
			`Here, together for the first time, are the shorter from Jim Butcher's DRESDEN FILES series — a compendium of cases that Harry and his cadre of allies managed to close in record time. The tales range from the deadly serious to the absurdly hilarious. Also included is a new, never-before-published novella that takes place after the cliff-hanger ending of the new April 2010 hardcover, *Changes*.

Contains:  
+ "Restoration of Faith"  
+ "Vignette"  
+ "Something Borrowed" — from * My Big Fat Supernatural Wedding*  
+ "It's My Birthday Too" — from * Many Bloody Returns*  
+ "Heorot" — from * My Big Fat Supernatural Honeymoon*  
+ "Day Off" — from * Blood Lite*  
+ "Backup" — novelette from Thomas' point of view, originally published by Subterranean Press  
+ "The Warrior" — novelette from * Mean Streets*  
+ "Last Call" — from * Strange Brew*  
+ "Love Hurts" — from * Songs of Love and Death*  
+ *Aftermath* — all-new novella from Murphy's point of view, set forty-five minutes after the end of * Changes*`},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			assert.Equal(t, td.x, HTML2Markdown(td.in), "unexpected markdown")
		})
	}
}

func TestHTML2Text(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name, in, x string
	}{
		{"empty string", "", ""},
		{"four dots", "....", "…"},
		{"four dots with spaces", ". . . .", "…"},
		{"three dots with spaces", ". . .", "…"},
		{"three dots", "...", "…"},
		{"dresden", `<b>HARRY DRESDEN — WIZARD</b><br /><br /><i>Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.</i><br /><br />Harry Dresden is the best at what he does. Well, technically, he's the <i>only</i> at what he does. So when the Chicago P.D. has a case that transcends mortal creativity or capability, they come to him for answers. For the "everyday" world is actually full of strange and magical things—and most don't play well with humans. That's where Harry comes in. Takes a wizard to catch a—well, whatever. There's just one problem. Business, to put it mildly, stinks.<br /><br />So when the police bring him in to consult on a grisly double murder committed with black magic, Harry's seeing dollar signs. But where there's black magic, there's a black mage behind it. And now that mage knows Harry's name. And that's when things start to get interesting.<br /><br />Magic - it can get a guy killed.`,
			`HARRY DRESDEN — WIZARD

Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.

Harry Dresden is the best at what he does. Well, technically, he's the only at what he does. So when the Chicago P.D. has a case that transcends mortal creativity or capability, they come to him for answers. For the "everyday" world is actually full of strange and magical things—and most don't play well with humans. That's where Harry comes in. Takes a wizard to catch a—well, whatever. There's just one problem. Business, to put it mildly, stinks.

So when the police bring him in to consult on a grisly double murder committed with black magic, Harry's seeing dollar signs. But where there's black magic, there's a black mage behind it. And now that mage knows Harry's name. And that's when things start to get interesting.

Magic - it can get a guy killed.`},
		{"dresden with link", `<i>An alternative cover edition with a different page count exists <a href="https://www.goodreads.com/book/show/13511897.here" title="here" rel="nofollow">here</a>.</i><br /><br />Harry Dresden - Wizard<br />Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.<br /><br />Harry Dresden has faced some pretty terrifying foes during his career. Giant scorpions. Oversexed vampires. Psychotic werewolves. It comes with the territory when you're the only professional wizard in the Chicago-area phone book.<br /><br />But in all Harry's years of supernatural sleuthing, he's never faced anything like this: The spirit world has gone postal. All over Chicago, ghosts are causing trouble - and not just of the door-slamming, boo-shouting variety. These ghosts are tormented, violent, and deadly. Someone - or <i>something</i> - is purposely stirring them up to wreak unearthly havoc. But why? And why do so many of the victims have ties to Harry? If Harry doesn't figure it out soon, he could wind up a ghost himself....`,
			`An alternative cover edition with a different page count exists here.

Harry Dresden - Wizard
Lost Items Found. Paranormal Investigations. Consulting. Advice. Reasonable Rates. No Love Potions, Endless Purses, or Other Entertainment.

Harry Dresden has faced some pretty terrifying foes during his career. Giant scorpions. Oversexed vampires. Psychotic werewolves. It comes with the territory when you're the only professional wizard in the Chicago-area phone book.

But in all Harry's years of supernatural sleuthing, he's never faced anything like this: The spirit world has gone postal. All over Chicago, ghosts are causing trouble - and not just of the door-slamming, boo-shouting variety. These ghosts are tormented, violent, and deadly. Someone - or something - is purposely stirring them up to wreak unearthly havoc. But why? And why do so many of the victims have ties to Harry? If Harry doesn't figure it out soon, he could wind up a ghost himself…`},
		{"dresden many links", `Here, together for the first time, are the shorter from Jim Butcher's DRESDEN FILES series — a compendium of cases that Harry and his cadre of allies managed to close in record time. The tales range from the deadly serious to the absurdly hilarious. Also included is a new, never-before-published novella that takes place after the cliff-hanger ending of the new April 2010 hardcover, <i>Changes</i>.<br /><br />Contains:<br />+ "Restoration of Faith"<br />+ "Vignette"<br />+ "Something Borrowed" -- from <em>
  <a href="https://www.goodreads.com/book/show/84156.My_Big_Fat_Supernatural_Wedding" title="My Big Fat Supernatural Wedding" rel="nofollow">My Big Fat Supernatural Wedding</a>
</em><br />+ "It's My Birthday Too" -- from <em>
  <a href="https://www.goodreads.com/book/show/140098.Many_Bloody_Returns" title="Many Bloody Returns" rel="nofollow">Many Bloody Returns</a>
</em><br />+ "Heorot" -- from <em>
  <a href="https://www.goodreads.com/book/show/1773616.My_Big_Fat_Supernatural_Honeymoon" title="My Big Fat Supernatural Honeymoon" rel="nofollow">My Big Fat Supernatural Honeymoon</a>
</em><br />+ "Day Off" -- from <em>
  <a href="https://www.goodreads.com/book/show/2871256.Blood_Lite" title="Blood Lite" rel="nofollow">Blood Lite</a>
</em><br />+ "<a href="https://www.goodreads.com/book/show/2575572.Backup" title="Backup" rel="nofollow">Backup</a>" -- novelette from Thomas' point of view, originally published by Subterranean Press<br />+ "The Warrior" -- novelette from <em>
  <a href="https://www.goodreads.com/book/show/3475145.Mean_Streets" title="Mean Streets" rel="nofollow">Mean Streets</a>
</em><br />+ "Last Call" -- from <em>
  <a href="https://www.goodreads.com/book/show/6122181.Strange_Brew" title="Strange Brew" rel="nofollow">Strange Brew</a>
</em><br />+ "Love Hurts" -- from <em>
  <a href="https://www.goodreads.com/book/show/7841656.Songs_of_Love_and_Death" title="Songs of Love and Death" rel="nofollow">Songs of Love and Death</a>
</em><br />+ <em>Aftermath</em> -- all-new novella from Murphy's point of view, set forty-five minutes after the end of <em>
  <a href="https://www.goodreads.com/book/show/6585201.Changes" title="Changes" rel="nofollow">Changes</a>
</em>`,
			`Here, together for the first time, are the shorter from Jim Butcher's DRESDEN FILES series — a compendium of cases that Harry and his cadre of allies managed to close in record time. The tales range from the deadly serious to the absurdly hilarious. Also included is a new, never-before-published novella that takes place after the cliff-hanger ending of the new April 2010 hardcover, Changes.

Contains:
+ "Restoration of Faith"
+ "Vignette"
+ "Something Borrowed" — from My Big Fat Supernatural Wedding
+ "It's My Birthday Too" — from Many Bloody Returns
+ "Heorot" — from My Big Fat Supernatural Honeymoon
+ "Day Off" — from Blood Lite
+ "Backup" — novelette from Thomas' point of view, originally published by Subterranean Press
+ "The Warrior" — novelette from Mean Streets
+ "Last Call" — from Strange Brew
+ "Love Hurts" — from Songs of Love and Death
+ Aftermath — all-new novella from Murphy's point of view, set forty-five minutes after the end of Changes`},
	}

	for _, td := range tests {
		td := td
		t.Run(td.name, func(t *testing.T) {
			assert.Equal(t, td.x, HTML2Text(td.in), "unexpected text")
		})
	}
}
