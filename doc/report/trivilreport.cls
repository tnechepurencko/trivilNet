\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{trivilreport}[ Trivil language report class]

\LoadClass[a4paper]{article}

\RequirePackage{xltxtra}
\RequirePackage{etoolbox}

\RequirePackage{polyglossia}
\setmainlanguage{russian}
\setotherlanguage{english}
\setkeys{russian}{babelshorthands=true}

\setmainfont{Times New Roman}
\setromanfont{Times New Roman} 
\setsansfont{Arial} 
\setmonofont{Courier New} 

\newfontfamily{\cyrillicfont}{Times New Roman} 
\newfontfamily{\cyrillicfontrm}{Times New Roman}
\newfontfamily{\cyrillicfonttt}{Courier New}
\newfontfamily{\cyrillicfontsf}{Arial}

\addto\captionsrussian{%
  \renewcommand{\figurename}{Рис.}%
  \renewcommand{\tablename}{Табл.}%
}

%
% Commands for language name and company name
%
\newcommand{\thelang}{Тривиль}
\newcommand{\theauthor}{Алексей Недоря}

%=============

\RequirePackage{iftex}

%\RequirePackage{unicode-math}
%\defaultfontfeatures{Scale=MatchLowercase}
%\defaultfontfeatures[\rmfamily]{Ligatures=TeX,Scale=1}

\RequirePackage{xcolor}

% Index commands
\RequirePackage{imakeidx}
\makeindex[columns=3, title=Alphabetical Index, intoc]

\RequirePackage{authblk}

\RequirePackage{listings}

% Highlights
\RequirePackage{soul}
\newcommand{\needlink}{\sethlcolor{yellow}\hl{link}}

\RequirePackage{hyperref}
\RequirePackage[numbered]{bookmark}

% Listings of grammar and source code

\lstset{% All listings settings
  backgroundcolor=\color{gray!15},
  numbers=left,
  numberstyle=\tiny,
  stepnumber=1,
  numbersep=5pt
}

% Grammar block style
\lstdefinestyle{thegrammar}{
  basicstyle=\small\ttfamily\color{blue!90}, %50
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


%==========
\lstdefinelanguage{sts}{
  sensitive = true,
  keywords={abstract,модуль,
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


%==================

\hypersetup{
  pdftitle={\thelang{} Language Report},
  hidelinks,
  pdfcreator={\theauthor},
  pdfauthor={\theauthor},
  pdfsubject={\thelang{} language report},
  pdfkeywords={\theauthor{}, \thelang{}, language report}
}

