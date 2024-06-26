# ASWSG - Another|Alexanders Static WebSite Generator

ASWSG allows you to generate Websites using Markup Syntax and HTML.

* *make friendly* - to be integrated in your workflow
* *configurable syntax* - to match different markup dialects

It is build with the idea of classical UNIX tools to do one job.
ASWSG behaves just like a compiler parsing a Markup file and generating a new HTML output.
ASWSG will not generate a filesystem structure, or an HTML frame, but you can do so using includes and
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

Note: This is under development, but what's not.

## Why should I use it

You like to create websites, and you want to have the freedom to build and use
html the way you like it using your favorite tool chain.
You are not forced to use a specific framework or project structure.

## Features

### Released

* [X] Simple markup parsing
  * [x] line based
    * [X] headers
    * [X] paragraphs
    * [X] ordered and unordered lists
    * [X] cites
    * [X] raw lines
    * [X] tables
    * [X] comments
  * [X] inline based
    * [X] bold
    * [X] emphasised
    * [X] strike trough
    * [X] underline
    * [X] links
      * [X] link type ``[[name|link]]`` (without rel=external)
      * [X] link type ``[name](link)`` (for external links with rel="external")
         * support link relationship attribute "external", see https://www.w3schools.com/tags/att_a_rel.asp
      * [x] link type ``[[URL]]``  with link-index
      * [X] link type ``=> URL description`` - a gemini style link
    * [X] block based
      * [X] quotes
      * [X] preformatted (code)
      * [X] raw html
      * [x] comments
* [X] include files
  * [X] include list of files by pattern
* [X] Setting and using of variables
  * [X] Setting variables (in this order, later overwrites former)
    1. static defaults
    2. dynamically like filename
    3. from command line parameter
    4. in parsed file
  * [x] Using variables everywhere, also in commands, but not in raw imported files
* [X] Commands
  * [X] dump-vars  (to log)
  * [X] message  (to log)
  * [X] include raw or crude files (with variable parsing and replacing)
  * [X] Execute an external command and insert its output
* [X] Conditional parsing based on variables
* [X] redefine markup tags
* [x] Log Filter
* [X] Anchor command (for headers)
* [X] Import CSV-file as table

## Usage

`aswsg [IN-FILE:]sourcefile [VARIABLE:value] > file.html`

Checkout the example.

## Control

### Messages

Which messages will be logged, respectively not logged, can be controlled through the
variable ASWSG-MESSAGE-FILTER.
Messages with severity in ASWSG-MESSAGE-FILTER will not be sent to *stderr*.


## Formatting

Description of the markup formatting.

You can use the mentioned variables to redefine the characters used for a formatting.

### Line level formatting

Used at begin the beginning line, using one of the characters.
Some characters are can be cascaded.

| Function                                                                                    | Default Char              | Variable                | Example                                                                                       |
|---------------------------------------------------------------------------------------------|---------------------------|-------------------------|-----------------------------------------------------------------------------------------------|
| Defining a variable                                                                         | ```@```                   | ```ASWSG-DEFINE```      | ```@variablename:value```                                                                     |
| Include a file oder files by pattern                                                        | ```+```                   | ```ASWSG-INCLUDE```     | ```+filename``` or ```+file_patern```                                                         |
| Raw (html) line to be inserted                                                              | ```$```                   | ```ASWSG-RAWLINE```     | ```$<article>```                                                                              |
| Escape for a paragraph char                                                                 | ```\```                   | ```ASWSG-ESCAPE```      | ```\* this is no bullet list```                                                               |
| Paragraph                                                                                   | (none)                    | (none)                  | ```Any text not starting not with a line level special. Empty lines start a new paragraph.``` |
| Header. The number of header characters define the depth of the header.                     | ```=``` or ```!```        | ```ASWSG-HEADER```      | ```== header level 2```                                                                       |
| Bullet list                                                                                 | ```*``` or ```-```        | ```ASWSG-LIST```        | ```* Bulltes and numbered Lists may be nested.```                                             |
| Numbered list                                                                               | any off ```#0123456789``` | ```ASWSG-NUMERATION```  | ```2# a level 2 indented list element```                                                      |
| Cite                                                                                        | ```>```                   | ```ASWSG-CITE```        | ```> To be or not to be.```                                                                   |
| Single line command, optionally closed by an ")", should not be changed                     | ```(```                   | ```ASWSG-COMMAND```     | ```(command parameter ...)```                                                                 |
| Defining a Table. The table character starts a new cell. One trailing pipe will be removed. | ```\|```                  | ```ASWSG-TABLE```       | ```\|a 2 cell\|table```                                                                       |
| Gemini style link                                                                           | ```=>```                  | ```ASWSG-GEMINI-LINK``` | ```=> http://haniarani.com/ Hania Rani```                                                     |
| Comment                                                                                     | ```;```                   | ```ASWSG-COMMENT```     | ```; this is a comment```                                                                     |

To continue a long line (i.e. a long header split over two lines) add an ```\``` to the end of the first line.


### Inline formatting

Used to format text within a line, using 2 or 3 strings for begin, end and middle when needed.

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

#### Glitches

There is a glitch when using the default markdown for emphasised ```//``` and a HTML link,
which usually contains also two dashes, in one line.
There is a small fix to avoid, that // will be used as the start of emphasised text,
when there is a ```:``` in front of it, but that prevents any later rendering of emphasised in this line.

A good and easy solution in this case is to redefine emphasised.
This could be done like this, to use two $-signs instead:

    @ASWSG-EMP-1:$$
    @ASWSG-EMP-2:$$

### Multi-Line/Block formatting

A line just containing at least three characters to enter a special block.
The block ends with the same characters in a line or a new block formatting.

| Function                        | Default Char | Variable            | Example    |
|---------------------------------|--------------|---------------------|------------|
| Citation                        | ```>```      | ```ASWSG-ML-CITE``` | ```>>>```  |
| Raw Lines                       | ```$```      | ```ASWSG-ML-RAW```  | ```$$$```  |
| Preformatted (code)             | ```%```      | ```ASWSG-ML-CODE``` | ```%%%```  |
| Comment                         | ```;```      | ```ASWSG-ML-COM```  | ```;;;```  |
| Horizontal line (just one line) | ```-```      | ```ASWSG-LINE```    | ```----``` |


## Commands

All commands start with the ```ASWSG-COMMAND``` character which default is ```(```.

| Function                                                                    | Form                                     |
|-----------------------------------------------------------------------------|------------------------------------------|
| Comment that will not be in the output HTML file                            | ```COMMENT any text```                   |
| Dump variables to log                                                       | ```DUMP-VARS parameters ignored```       |
| Write a message to the log                                                  | ```MESSAGE any text```                   |
| Set a link anchor with ```name```                                           | ```ANCHOR anchor-name```                 |
| Insert link index for Link-3 of write it to ```[<file>]```                  | ```LINK-INDEX filename```                |
| Include file raw without variable substitution                              | ```INCLUDE-FILE-RAW filename```          |
| Include file crude without variable substitution                            | ```INCLUDE-FILE-CRUDE filename```        |
| Include the output of a script/command (no variable substitution)           | ```INCLUDE-SCRIPT programm parameters``` |
| Include CSV-file as table                                                   | ```INCLUDE-CSV filename.csv```           |
| Use the line formatter ```>``` as ```<blockquote>``` instead off ```<cite>```. Activate this with the parameter "TRUE", "T" or no parameter. Every other parameter will be interpreted as false. (Default is ```cite```.) | ```GT-AS-BLOCKQUOTE``` |

## Conditional commands

If a condition is set the following lines are only parsed if the condition criteria are true.
This is valid until a new condition is set, a ```COND-ELSE``` or ```COND-END``` clause changes this.
There is only one condition valid at a time.
There is no nesting, a new ```COND-IF-...``` just changes the condition.

All conditions start with the ```ASWSG-COMMAND``` character which default is ```(```.

| Function                                      | Form                                           |
|-----------------------------------------------|------------------------------------------------|
| Check if Variable is set                      | ```COND-IF-SET variable-name```                |
| Check if Variable is not set                  | ```COND-IF-NOT-SET variable-name```            |
| Check if Variable has the value               | ```COND-IF-EQUAL variable-name value...```     |
| Check if Variable has not the value           | ```COND-IF-NOT-EQUAL variable-name value...``` |
| Condition else case (last condition inverted) | ```COND-ELSE```                                |
| End condition validation                      | ```COND-END```                                 |

## More Variables

### Control Variables

| Function                                                                                                                   | Name:Default                                         |
|----------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------|
| Number of header lines when parsing a table.                                                                               | ```ASWSG-TABLE-HEADERLINES:1```                      |
| Alignment of colons when parsing a table. L=left, C=center, R=right. If the value is to short, or unkown it defaults to L. | ```ASWSG-TABLE-ALIGNMENT:LL```                       |
| CSV field separator                                                                                                        | ```ASWSG-CSV-COMMA:;```                              |
| CSV comment line                                                                                                           | ```ASWSG-CSV-COMMENT:#```                            |
| Automatically generate anchors for headers, T = true, everything else is false                                             | ```ASWSG-AUTO-GENERATE-ANCHOR:T```                   |
| Reverse Order of included files, T = true, everything else is false                                                        | ```ASWSG-INCLUDE-REVERSE:F```                        |
| Date format                                                                                                                | ```DATEFORMAT:2006-01-02```                          |
| Time format                                                                                                                | ```TIMEFORMAT:15:04:05```                            |
| Timestamp format                                                                                                           | ```TIMESTAMPFORMAT:2006-01-02 15:04:05 UTC+ 07:00``` |
| Filter out message types, default is "Dd" for debug messages                                                               | ```ASWSG-MESSAGE-FILTER:Dd```                        |
| Use Gemini style links ```=>```,T = true, everything else is false                                                         | ```GEMINI-LINK```                        |

The date and time format refers to GO's [Time.Format](https://pkg.go.dev/time#pkg-constants).

### Info Variables

| Function                          | Name                |
|-----------------------------------|---------------------|
| Name of the file currently parsed | ```FILENAME```      |
| Name of the main file             | ```IN-FILE```       |
| Current date                      | ```TODAY```         |
| Current time                      | ```TIME```          |
| Current timestamp (date + time)   | ```NOW```           |
| The ASWSG version                 | ```ASWSG-VERSION``` |
| The ASWSG authors name            | ```ASWSG-AUTHOR```  |
| The ASWSG license                 | ```ASWSG-LICENSE``` |


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

    $<! a raw line></body></html>


## Parsing Logic

* Start (main)
   * setDefaultSiteVars
   * parseAndSetCommandLineVars
   * parseFile, for each line
       * process continued lines
       * parseLine
            * replaceInlineVars
            * parseCondition
            * validateCondition (return if not fulfilled)
            * block mode comment (return if true)
            * block mode code (return if true)
            * block mode crude (return if true)
            * block mode cite (continue) (TODO)
            * empty Line (new paragraph state)
            * ignore comment line
            * parseAndSetVar (return if true)
            * parseCommands (return if true)
            * parse raw lines (return if true)
            * process includes (parseFile, return if true)
            * parse gemini link (return if true)
            * parse header (return if true)
            * parseTableLine
            * parse horizontal line
            * parseCommonParagraphControls and regular text lines
                * parse LIST, CITE and NUMERATION
                * parse escaped line
                * parse paragraph / empty line
                * surroundWithHTMLTag for list, cite, numeration
                * parseInLine for
                    * bold, emphasised, strike, code, link(1-3)


