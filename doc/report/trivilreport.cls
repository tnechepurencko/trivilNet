\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{trivilreport}[2023/01/16 Trivil language report class]

\LoadClass[a4paper]{article}

%
% Commands for language and author 
%
\newcommand{\thelang}{Тривиль}
\newcommand{\theauthor}{Алексей Недоря}

\RequirePackage[hyphens]{url}

\RequirePackage{amsmath}
\RequirePackage{amssymb}
\RequirePackage{lmodern}
\RequirePackage{iftex}

\RequirePackage{unicode-math}
\defaultfontfeatures{Scale=MatchLowercase}
\defaultfontfeatures[\rmfamily]{Ligatures=TeX,Scale=1}

% Use upquote if available, for straight quotes in verbatim environment
\IfFileExists{upquote.sty}{
  \RequirePackage{upquote}
}{}

\IfFileExists{microtype.sty}{
  \RequirePackage[]{microtype}
  \UseMicrotypeSet[protrusion]{basicmath}
}{}
             
\makeatletter
\@ifundefined{KOMAClassName}{% if non-KOMA class
  \IfFileExists{parskip.sty}{%
    \RequirePackage{parskip}
  }{% else
    \setlength{\parindent}{0pt}
    \setlength{\parskip}{6pt plus 2pt minus 1pt}}
}{% if KOMA class
  \KOMAoptions{parskip=half}}
\makeatother

\RequirePackage{xcolor}
\RequirePackage{longtable}
\RequirePackage{booktabs}
\RequirePackage{array}

% For calculating minipage widths
\RequirePackage{calc} 

% Correct order of tables after \paragraph or \subparagraph
\RequirePackage{etoolbox}

\makeatletter
\patchcmd\longtable{\par}{\if@noskipsec\mbox{}\fi\par}{}{}
\makeatother

% Allow footnotes in longtable head/foot
\IfFileExists{footnotehyper.sty}{
  \RequirePackage{footnotehyper}
}{
  \RequirePackage{footnote}
}

\makesavenoteenv{longtable}

\RequirePackage[normalem]{ulem}
\setlength{\emergencystretch}{3em} % prevent overfull lines

\IfFileExists{bookmark.sty}{
  \RequirePackage[numbered]{bookmark}
}{
}

% Add URL line breaks if available
\IfFileExists{xurl.sty}{
  \RequirePackage{xurl}
}{}
\urlstyle{same} % disable monospaced font for URLs

% Index commands
\RequirePackage{imakeidx}
\makeindex[columns=3, title=Alphabetical Index, intoc]

\RequirePackage{authblk}

% Listings of grammar and source code
\RequirePackage{listings}

\lstset{% All listings settings
%  backgroundcolor=\color{gray!15},
  numbers=left,
  numberstyle=\tiny,
  stepnumber=1,
  numbersep=5pt
}

% Grammar block style
\lstdefinestyle{thegrammar}{
  basicstyle=\small\ttfamily\color{blue!50},
  breaklines=true,
  moredelim=[s][\color{green!50!black}\ttfamily]{'}{'},
  moredelim=*[s][\color{black}\ttfamily]{options}{\}},
  commentstyle={\color{gray}\itshape},
  morecomment=[l]{//},
  emph={%
    STRING,
  },
  emphstyle={\color{blue}\ttfamily},
  alsoletter={:,|,;},
  morekeywords={:,|,;},
  keywordstyle={\color{black}},
}

\lstnewenvironment{grammar}{
  \lstset{style=thegrammar}
}{}

% Internal links
\RequirePackage{xparse}

\hypersetup{
  pdftitle={Язык программирования \thelang{}},
  hidelinks,
  pdfcreator={\theauthor},
  pdfauthor={\theauthor},
  pdfsubject={Описания языка программирования \thelang{} },
  pdfkeywords={\theauthor{}, \thelang{}, описание языка программирования}
}

\makeatletter

% hyperref should be last import
\RequirePackage{hyperref}

% Specification commands

% TODO Change name, add star version
%\newcommand{\See}[1]{\textit{see} \S\ref{#1}\seechecknextarg}
\newcommand{\See}[1]{\S\ref{#1}\seechecknextarg}
\newcommand{\seechecknextarg}{\@ifnextchar\bgroup{\seegobblenextarg}{}}
\newcommand{\seegobblenextarg}[1]{\@ifnextchar\bgroup{, \S\ref{#1}\seegobblenextarg}{ \textit{and} \S\ref{#1}}}

\newcommand{\Seename}[1]{
  \see ``\nameref{#1}''\seenamechecknextarg}
\newcommand{\seenamechecknextarg}{\@ifnextchar\bgroup{\seenamegobblenextarg}{}}
\newcommand{\seenamegobblenextarg}[1]{\@ifnextchar\bgroup{, ``\nameref{#1}''\seenamegobblenextarg}{ and ``\nameref{#1}''}}

\makeatother

\newcommand{\code}[1]{\texttt{#1}}

\NewDocumentCommand \definition { s m }{
  \IfBooleanTF{#1}
              {\index{#2|hyperpage}\emph{#2}}
              {\index{\lowercase{#2}|hyperpage}\emph{#2}}
}

\NewDocumentCommand \mention { s m }{
  \IfBooleanTF{#1}
              {\index{#2|hyperpage}#2}
              {\index{\lowercase{#2}|hyperpage}#2}
}

% Highlights
\RequirePackage{soul}
\newcommand{\notsolved}[1]{\sethlcolor{magenta}\hl{#1}}
\newcommand{\needlink}[1]{\sethlcolor{yellow}\hl{#1}}
\newcommand{\notwritten}[1]{{\color{white}\sethlcolor{gray}\hl{#1}}}
\newcommand{\tobediscussed}[1]{\sethlcolor{cyan}\hl{#1}}
\newcommand{\tobedeleted}[1]{\sethlcolor{orange}\hl{#1}}

% Source highlight


\lstdefinelanguage{sts}{
  sensitive = true,
  keywords={abstract,as,assert,await,break,case,
    class,const,constructor,continue,defer,do,
            else,enum,export,false,for,function,if,import,inner,instanceof,
            interface,launch,let,native,new,null,open,override,package,private,protected,
            public,return,static,switch,super,this,throw,true,try,type,while}
  otherkeywords={% Operators
    +,-,*,/,\%,&,|,^,\<\<,\>\>,+=,-=,*=,/=,\%=,&=,|=,^=,
    \<\<=,\>\>=,\>\>\>=,\&\&,||,++,--,==,<,>,==,==,!,??,
    ?.,!=,<=,>=,...,\(,\[,\{,\,,\.,\),\],\},;,:
  },
  basicstyle=\ttfamily\small,
  keywords = [2]{catch,default,extends,from,implements,in,of,out,recover,rethrows,throws},
  keywordstyle=\color{blue},
  keywordstyle=[2]\color{blue},% for example
  numbers=left,
  numberstyle=\scriptsize,
  stepnumber=1,
  numbersep=8pt,
  showstringspaces=false,
  breaklines=true,
  frame={single},
  comment=[l]{//},
  morecomment=[s]{/*}{*/},
  commentstyle=\color{purple}\ttfamily,
  stringstyle=\color{red}\ttfamily,
  morestring=[b]',
  morestring=[b]"
}

\lstnewenvironment{source}{
  \lstset{language=sts}
}{}

% Grammar block style
\lstdefinestyle{theconsole}{
  basicstyle=\ttfamily\color{black!70},
  stepnumber=1,
  numbersep=8pt,
  showspaces=true,
  breaklines=true,
  frame={single},
}

\lstnewenvironment{console-output}{
  \lstset{style=theconsole}
}{}
