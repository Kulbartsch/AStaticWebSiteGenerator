# ASWSG - Another|Alexanders Static WebSite Generator

ASWSG allows you to generate Websites using Markup Syntax and HTML.
It is build with the idea of classical unix tools to do one job.
ASWSG behaves just like a compiler parsing a Markup file and generating a new HTML output.
ASWSG will not generate a file structure or an HTML frame, but you can do so using includes and
use it with a build system like *make*.

**Please note: This is under development. The output code must currently be manually be cleaned before use.**

HTML and Markup can be mixed, reusable includes may be used for headers, footers and other repeating code blocks.
Most benefits come from using dynamic variables, which are partly set dynamically by aswsg,
by the document or via the command line and the include system.
In example the current *date*, *time* and the actual *filename* is set dynamically, *author*, *creation date*
and others in the document, the *sitename* as a command line parameter.
But it is even possible to redefine markup symbols due to different markup dialects.
(i.E. use "!" instead/as well of "=" for headers.)
All variables can be used subsequent, also in included files. (i.E. using the article name in a header include.)

This tool will generate new HTML code, which -- of course -- may contain dynamic code.

## (Planned) Features

may vary

### Release 1

* [X] Simple markup parsing
  * [x] line based
    * [X] headers
    * [X] paragraphs
    * [X] unordered lists (~~1 level~~ nested)
    * [x] ordered lists (~~1 level~~ nested)
    * [X] cites
    * [X] raw lines
  * [X] inline based
    * [X] bold
    * [X] emphasised
    * [X] strike trough
    * [X] links
    * [X] code
* [X] include parsed files
* [X] Setting and using of variables
  * [X] Setting variables (in this order, later overwrites former)
    1. [X] static
    2. [X] dynamically
    3. [X] from command line parameter
    4. [X] in file
  * [x] Using
* [X] Commands
  * [X] dump-vars  (to log)
  * [X] message  (to log)
* [X] redefine markup tags
* [X] *make* friendly
* [X] go tests (partly)
* [X] An example

### Future Releases

* Commands
  * [ ] comments
  * [ ] dump-context  (to log)
  * [ ] interactive  (enter interactive mode = read from io.stdin)
  * [ ] execute-shell-command  <command with parameters>
  * [ ] include raw files
  * [ ] include raw files, but with with variable parsing and replacing
  * [ ] execute-script <filename>  (run a script ... maybe in the future)
* [ ] Log Filter
* [ ] Inherit of HTML code (without using raw command)
* [ ] Simple markup parsing
  * [ ] multi line
    * [ ] raw
    * [ ] cite
    * [ ] code
* [ ] Anchor for headers
* [ ] more examples
* [ ] more tests
* [ ] Index of page (based on header)
* [ ] More link types
  * [ ] camelCase links
  * [ ] automatic URL detection
  * [ ] ``[[ ]]`` links without URL, auto generates internal links
  * [ ] link type ``[name](link)``
* [ ] individual HTML tag IDs and classes
* [ ] increased markup features like
  * [ ] Tables
  * [ ] Pictures / embeded documents
* [x] nested lists
* [ ] multi line lists
* [ ] Basic markdown features

## Usage

Description of the markup formating.

### Line level formating

Used at begin the beginning line, using one of the characters.
Some characters are can be cascaded.

#### Defining a variable

Variable: "ASWSG-DEFINE"

Default character: "@",

Special: Define a variable.

Format: @variablename:value

#### Include a file

Variable: "ASWSG-INCLUDE"

Default character: "+"

special: include a text

Format: +filename

#### raw (html) line

Variable: "ASWSG-RAWLINE"

Default character: "$"

Special: line will be inserted as is

Format: ```$<article>```


#### Escape

Variable: "ASWSG-ESCAPE"

Default character: "\"

Special: special: escape char for paragraph

Format: ```\* this is no bullet list```


#### paragraph

Variable: (none)

Default character: (none)

Special: (none)

Format: any text not starting not with a line level special.

Empty lines start a new paragraph.


#### Bullet list

Variable: "ASWSG-LIST"

Default characters: "*" and "-"

Format:

    * This is
    * just a simple List
    *- with four entries
    *- in two levels

Bulltes and numbered Lists may be nested.

#### Cite

Variable: "ASWSG-CITE"

Default character: ">"

Format: > To be or not to be.

#### Numbered list

Variable: "ASWSG-NUMERATION"

Default character: "#0123456789"

Format:

   	# A numbered list
   	1 can be made with the numbers 0-9
   	1 for your convenience
   	21 but if you use two digits
   	22 it will be handeled like a nested list
   	5## numbers and # sign can be mixed
   	1 numbers don't have to be in sequence

#### Commands (not implemented yet)

Variable: "ASWSG-COMMAND"

Default character: "("

Special: single line command, optionally closed by an ")", should not be changed

Format: ```(command)```


#### Defining a Table (not implemented yet)

Variable: "ASWSG-TABLE"

Default character: "|"

Special:

Format:


#### Header

Variable: "ASWSG-HEADER"

Default character: "=!"

Special: number of header characters define the depth of the header

Format: ```== header level 2```
