# ASWSG - Another|Alexanders Static WebSite Generator

**Please note: This is the first, currently unusable, code.**

Under development.

ASWG allows you to generate Websites using HTML, and Markup Syntax. It is useable with a build system like make.

HTML and Markup can be mixed, reusable
Most benefits come from using dynamic variables, which are partly set dynamically by aswsg, by the document or via the command line and the include system.
In example the current *date*, *time* and the actual *filename* is set dynamically, *author*, *creation date* and others in the document, the *sitename* as a command line parameter.
But it is even possible to redefine markup symbols due to different markup dialects.
(i.E. use "!" instead/as well of "=" for headers.)
All variables can be used subsequent, also in included files. (i.E. using the article name in a header include.)

This tool will generate new HTML code, which -- of course -- may contain dynamic code.

## Planned Features

may vary

### Release 1

* [ ] Simple markup parsing
  * [ ] line based
    * [ ] headers (currently only 1 level)
    * [X] paragraphs
    * [X] unordered lists (~~1 level~~ nested)
    * [x] ordered lists (~~1 level~~ nested)
    * [ ] cites
    * [ ] raw lines
  * [ ] inline based
    * [X] bold
    * [ ] emphasised
    * [ ] strike trough
    * [ ] links
    * [ ] code
  * [ ] multi line
    * [ ] raw
    * [ ] cite
    * [ ] code
* [X] Include file processing
* [X] Setting and using of variables
  * [X] Setting variables (in this order, later overwrites former)
    1. [X] static
    2. [X] dynamically
    3. [X] from command line parameter
    4. [X] in file
  * [x] Using
* [x] redefine markup tags
* [x] *make* friendly
* [ ] go tests (partly)
* [ ] An example

### Future Releases

* [ ] Anchor for headers
* [ ] more examples
* [ ] Index of page (based on header)
* [ ] More link types
  * [ ] camelCase links
  * [ ] automatic URL detection
  * [ ] more features for [[ ]] links
* [ ] individual HTML tag IDs and classes
* [ ] increased markup features like
  * [ ] Tables
  * [ ] Pictures / embeded documents
* [x] nested lists
* [ ] multi line lists
* [ ] Basic markdown features
