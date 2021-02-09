# goforit
A pure go package for enabling runtime customization formulas for Go utilizing JavaScript package. 

With Go-For-It package, developers can externalize program formula(s) (e.g., in text file, csv, excel file, DB, etc.) and load them during runtime. 

The benefit of having formula(s) externalized is, for program that extensively uses formula(s) so numbers of formula(s) can be maintained outside of source code which in turn making the source code more cleaner, adding/changing formula can be done without stopping the running program, as well as complex formula(s) could be handled by specialists (a.k.a. users).  

**NOTE** The current implementation is built on top of "otto"
