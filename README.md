# expornote

Expornote is a quick-n-dirty utility I wrote for extracting notes and
attachments from Evernote's `.enex` export files so that I could move them to
Dropbox.

## Output Format

* All attachments and notes are extracted to a new directory alongside your
  `.enex` file. If your export is named "MyNotes.enex", the directory will be
  named "MyNotes".

* Attachments are exported as `<note_title>.<extension>`.

* Multiple attachments are exported as `<note_title> (1).<extension>` (and so
  on).

* Note contents are converted to Markdown as cleanly as possible. This is not
  perfect, however, so I suggest keeping a backup of your `.enex` files if you
  can't verify the contents of all of your exported notes.

## Usage

Install `pandoc` and `expornote`:

    brew install pandoc
    go install github.com/tysonmote/expornote

Run `expornote` on your exported Evernote notes:

    expornote MyNotes.enex

## Limitations

* I don't use Evernote for notes (just PDF storage, mostly), so I didn't tune
  the Markdown much. It seems to be reasonable, but Pull Requests are welcome!

* I've only tested `expornote` against my Evernote exports, so it's possible
  that there's bugs. Again, Pull Requests are welcome!
