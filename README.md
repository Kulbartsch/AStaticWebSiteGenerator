# ASWSG - Another|Alexanders Static WebSite Generator

ASWSG allows you to generate Websites using Markup Syntax and HTML.

* *make friendly* - to be integrated in your workflow
* *adaptable syntax* - to match different markup dialects

It is build with the idea of classical UNIX tools to do one job.
ASWSG behaves just like a compiler parsing a Markup file and generating a new HTML output.
ASWSG will not generate a file structure or an HTML frame, but you can do so using includes and
use it with a build system like *make*.

HTML and Markup can be mixed, reusable includes may be used for headers, footers and other repeating code blocks.
Most benefits come from using dynamic variables, which are partly set dynamically by aswsg,
by the document or via the command line and the include system.
In example the current *date*, *time* and the actual *filename* is set dynamically, *author*, *creation date*
and others in the document, the *sitename* as a command line parameter.
But it is even possible to redefine markup symbols due to different markup dialects.
(i.E. use "!" instead/as well of "=" for headers.)
All variables can be used subsequent, also in included files. (i.E. using the article name in a header include.)

This tool will generate new HTML code, which -- of course -- may contain dynamic code.

Note: This is under development but usable.

## Features

### Released

* [X] Simple markup parsing
  * [x] line based
    * [X] headers
    * [X] paragraphs
    * [X] unordered lists (~~1 level~~ nested)
    * [x] ordered lists (~~1 level~~ nested)
    * [X] cites
    * [X] raw lines
    * [X] tables
  * [X] inline based
    * [X] bold
    * [X] emphasised
    * [X] strike trough
    * [X] links
      * [X] link type ``[[name|link]]`` (without rel=external)
      * [X] link type ``[name](link)`` (for external links with rel="external")
      * [X] support link relationship attribute "external", see https://www.w3schools.com/tags/att_a_rel.asp
      * [x] link type ``[[URL]]``  with link-index
    * [X] code
* [X] include parsed files
* [X] continued lines
* [X] Setting and using of variables
  * [X] Setting variables (in this order, later overwrites former)
    1. [X] static
    2. [X] dynamically
    3. [X] from command line parameter
    4. [X] in file
  * [x] Using variables everywhere, also in commands, but not in raw imported files
* [X] Commands
  * [X] dump-vars  (to log)
  * [X] message  (to log)
  * [X] comments
  * [X] dump-context  (to log)
  * [X] Anchor
  * [X] include raw files
  * [X] include crude files, but with with variable parsing and replacing
  * [X] Execute an external command and insert its output
* [X] redefine markup tags
* [X] *make* friendly
* [X] go tests (partly)
* [X] An example
* [X] continued lines (for long text)
* [x] Log Filter
* [X] Anchor for headers

### Future Releases

* [ ] Command to import list of files by pattern

### More Ideas

* [ ] Automatic conversion of html sensitive chars to html (<>&)
* [ ] Option to generate HTML boilerplate &lt;html>,&lt;body> and CSS
* [ ] Pictures / embeded documents (IMAGE &lt;filename> &lt;Alt Text>)

* [ ] raw blocks
* [ ] cite blocks
* [ ] code blocks

* [ ] individual HTML tag classes

* [ ] Index of page (based on header)

* [ ] Command to Execute an external command  &lt;command with parameters>
* [ ] Command to include CSV file as table
* [ ] camelCase links
* [ ] automatic URL detection
* [ ] HTML classes for microformats http://microformats.org/


## Usage

```aswsg [IN-FILE=]sourcefile [VARIABLE:value] > file.html```

Checkout the example.

## Controll

### Messages

Which messages will be logged, respectivly not logged, can be controlled thru the
variable ASWSG-MESSAGE-FILTER.
Messages with severity in ASWSG-MESSAGE-FILTER will not be send to *stderr*.


## Formating

Description of the markup formating.

You can use the mentioned variables to redefine the characters used for a formating.

### Line level formating

Used at begin the beginning line, using one of the characters.
Some characters are can be cascaded.

| Function | Default Char | Variable | Example |
| -------- | ------------ | -------- | ------- |
| Defining a variable | ```@``` | ```ASWSG-DEFINE``` | ```@variablename:value``` |
| Include a file | ```+``` | ```ASWSG-INCLUDE``` | ```+filename``` |
| Raw (html) line to be inserted | ```$``` | ```ASWSG-RAWLINE``` | ```$<article>``` |
| Escape for a paragraph char | ```\``` | ```ASWSG-ESCAPE``` | ```\* this is no bullet list``` |
| Paragraph | (none) | (none) | ```Any text not starting not with a line level special. Empty lines start a new paragraph.``` |
| Header. The number of header characters define the depth of the header. | ```=``` or ```!``` | ```ASWSG-HEADER``` | ```== header level 2``` |
| Bullet list | ```*``` or ```-``` | ```ASWSG-LIST``` | ```* Bulltes and numbered Lists may be nested.``` |
| Numbered list | any off ```#0123456789``` | ```ASWSG-NUMERATION``` | ```2# a level 2 indented list element``` |
| Cite | ```>``` | ```ASWSG-CITE``` | ```> To be or not to be.``` |
| Single line command, optionally closed by an ")", should not be changed | ```(``` | ```ASWSG-COMMAND``` | ```(command parameter ...)``` |
| Defining a Table. The table character starts a new cell. One trailing pipe will be removed. | ```\|``` | ```ASWSG-TABLE``` | ```\|a 2 cell\|table``` |

To continue a long line (i.e. a long header split over two lines) add an ```\``` to the end of the first line.


## Inline formating

Used to format text within a line line, using 2 or 3 strings for begin, end and middle when needed.

| Function | Default Char | Variables | Example |
| -------- | ------------ | -------- | ------- |
| Variable replacing |  ```{{```...```}}``` | ```ASWSG-VAR-1```...```ASWSG-VAR-2``` | ```{{```variablename```}}```  |
| Bold | ```*```...```*``` | ```ASWSG-BOLD-1```...```ASWSG-BOLD-2``` | ```*```bold text```*```  |
| Emphasised | ```//```...```//``` | ```ASWSG-EMP-1```...```ASWSG-EMP-2``` | ```//```italic text```//```  |
| Code | ``` `` ```...``` `` ``` | ```ASWSG-CODE-1```...```ASWSG-CODE-2``` | ``` `` ```text``` `` ```  |
| Strike through | ```~~```...```~~``` | ```ASWSG-STRIKE-1```...```ASWSG-STRIKE-2``` | ```~~```text```~~```  |
| Underline | ```__```...```__``` | ```ASWSG-UNDERL-1```...```ASWSG-UNDERL-2``` | ```__```text```__```  |
| Link-1 (internal) | ```[[```...```\|```...```]]``` | ```ASWSG-LINK-1-1```...```ASWSG-LINK-1-3```...```ASWSG-LINK-1-2``` | ```[[```text```\|```URL```]]```  |
| Link-2 (with rel=external) | ```[```...```](```...```)``` | ```ASWSG-LINK-2-1```...```ASWSG-LINK-2-3```...```ASWSG-LINK-2-2``` | ```[```text```](```URL```)```  |
| Link-3 (Link Index, to be inserted with commnd (LINK-INDEX)) | ```[[```...```]]``` | ```ASWSG-LINK-3-1```...```ASWSG-LINK-1-2``` | ```[[```URL```]]```  |


### Multi-Line/Block formating

A line just containing at least three characters to enter a special block.
The block ends with the same characters in a line or a new block formating.

| Function | Default Char | Variable | Example |
| -------- | ------------ | -------- | ------- |
| Citeation (tbi) | ```>``` | ```ASWSG-ML-CITE``` | ```>>>``` |
| Raw Lines (tbi) | ```$``` | ```ASWSG-ML-RAW``` | ```$$$``` |
| Code (tbi) | ```%``` | ```ASWSG-ML-CODE``` | ```%%%``` |
| Horizontal line (just one line) | ```-``` | ```ASWSG-LINE``` | ```----``` |

tbi = to be implemented, does not exist jet.

## Commands

| Function | Form |
| -------- | ------------ |
| Comment that will not be in the output HTML file | ```COMMENT any text``` |
| Dump variables to log | ```DUMP-VARS parameters ignored``` |
| Write a message to the log | ```MESSAGE any text``` |
| Set an link anchor | ```ANCHOR anchor-name``` |
| Insert link index  | ```LINK-INDEX parameters ignored``` |
| Include file raw without variable substitution | ```INCLUDE-FILE-RAW filename``` |
| Include file crude without variable substitution | ```INCLUDE-FILE-CRUDE filename``` | 

## More Variables

### Control Variables

| Function | Name:Default |
| -------- | ------------ |
| Number of header lines when parsing a table. | ```ASWSG-TABLE-HEADERLINES:1``` |
| Alignment of colons when parsing a table. L=left, C=center, R=right. If the value is to short, or unkown it defaults to L. | ```ASWSG-TABLE-HEADERLINES:LL``` |
| Automatically generate anchors for headers, T = true, everything else is false | ```ASWSG-AUTO-GENERATE-ANCHOR:T``` |
| Date format | ```DATEFORMAT:2006-01-02``` |
| Time format | ```TIMEFORMAT:15:04:05``` |
| Timestamp format | ```TIMESTAMPFORMAT:2006-01-02 15:04:05 UTC+ 07:00``` |

### Info Variables

| Function | Name |
| -------- | ------------ |
| Name of the file currently parsed | ```FILENAME``` |
| Name of the main file | ```IN-FILE``` |
| Current date | ```DATE``` |
| Current time | ```TIME``` |
| Current timestamp (date + time) | ```NOW``` |


## Example


    (comment next line defines a variable title)
    @title:Test Page
    
    +include filename

    == A Header {{title}}
    
    * This is
    * just a simple List
    *- with four entries
    *- in two levels
    
   	# A numbered list
   	1 can be made with the numbers 0-9
   	1 for your convenience
   	21 but if you use two digits
   	22 it will be handeled like a nested list
   	5## numbers and # sign can be mixed
   	1 numbers don't have to be in sequence

    | a table | with two columns
    | and two | rows

    $<! a raw line></bosy></html>
