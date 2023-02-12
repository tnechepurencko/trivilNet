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

%\RequirePackage{iftex}

%\RequirePackage{unicode-math}
%\defaultfontfeatures{Scale=MatchLowercase}
%\defaultfontfeatures[\rmfamily]{Ligatures=TeX,Scale=1}

\RequirePackage{xcolor}

% Index commands
\RequirePackage{imakeidx}
\makeindex[columns=3, title=Alphabetical Index, intoc]

%\RequirePackage{authblk}

\RequirePackage{fancyvrb}
\renewcommand{\theFancyVerbLine}{%
{\footnotesize\arabic{FancyVerbLine}}}

\DefineVerbatimEnvironment%
{Trivil}{Verbatim}
{frame=single,numbers=left}

\DefineVerbatimEnvironment%
{Grammar}{Verbatim}
{frame=lines,formatcom=\color{blue},rulecolor=\color{black}}

% Highlights
\RequirePackage{soul}
\newcommand{\needlink}{\sethlcolor{yellow}\hl{link}}

\RequirePackage{hyperref}
\RequirePackage[numbered]{bookmark}

%==================

\hypersetup{
  pdftitle={\thelang{} Language Report},
  hidelinks,
  pdfcreator={\theauthor},
  pdfauthor={\theauthor},
  pdfsubject={\thelang{} language report},
  pdfkeywords={\theauthor{}, \thelang{}, language report}
}

