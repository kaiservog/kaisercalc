## Kaisercalc purpose
    
Kaisercalc is a solution for resolving calculations that is a pain in the a** to resolve in simple calculators, it provides variable/function definition. It reads files with math expression and can print the result and can resolve math expression in one line CLI.

---
### Line types
there are three types of line to inform you calculations:
* Declare
* System functions
* Imports

**Declare**
It's define the value/function to a variable 
i.e
``PI=3.1415`` or ``months=12`` or setting a functions ``salary(x)=x*12``

**System functions**
There are only two system functions **print** and **println**, they both prints values to the screen.
* print - prints the value in the screen
* println - prints the value and jump to a new line in the screen

i.e
``print(2+2)`` or  ``println(X)`` or ``print("the result is", x)``

prints call can receive more than one elements to print and it can be anything that results in a number, like variables, functions call and a special type "string"
i.e
`print(1, 2, 3)` it prints `1 2 3 ` (a space after every element is added)
`print(PI)` it prints `3.1 ` (assuming PI is defined as 3.1)
`println(year_salary(1000), 2)` it prints 
`12000`
`2`
(assuming year_salary(x) is `x*12`)
`print('the result is', 42)` it prints `the result is 42`

so an example of a calculation file looks like
```
employees= 120
salary=2200
taxes(x)=x*0.1+20
println('employees #', employees, 'payment is', taxes(employees)+salary*employees)
```

**imports**
Some times you may want to split some data into files like `finance.txt`, `study.txt` or `test02.txt`
that way you may have the need to import data from another file, for example
you can have a file to declare some personal data like account ballance, bills info ... when calculating currently month bills value, you can create a file january_bill.txt and import the bills.txt data, now you can access variables and functions defined in the bills.txt, next month you can delete january file and create another one.

the syntax to use import is `import alias file_name.txt`
i.e
`import bills bills.txt` `import prob functions/probs.txt`
---
### Variables
Variables must stats with a letter than can has "_" (underscore) in name.
Can have capital letters or not
i.e
`PI=3` `pi=3.1` `num_of_books=14`

if the set of a variable is a functions it must inform the new variables used in function
i.e
`
calc_serie_time(minutes_per_episode, episodes, breaks)=minutes_per_episode*episodes + breaks*episodes*5
`

in this case the new variables used in function is `minutes_per_episode`, `episodes` and `breaks`; 
All of then must be inside of `( )`in the variable name  
variables can be defined as text as well, in that case we define them with `''` like `impact_phrase='i am a text'`
---
### System functions
Can only be used in lines without a definition, so this code is not correctly
i.e
`PI=print('what?')`
the correctly way to use is in the start of a line
i.e
`print('pi is', PI)`

print and println can receive any variable, number or text:
i.e
`print(1, 1+1, calc(x), 'hi!', PI)` it prints `1 2 3 hi! 5`
---
### Calling kaisercalc
**file**
Just use `kaisercalc filename`
i.e
`kaisercalc mycalc.txt` or in Windows `kaisercalc.exe mycalc.txt`

**simple expression**
You can use kaisercalc to resolve expression in CLI, in that case use the argument `-e` `kaisercalc -e "expression"`
i.e
`kaisercalc -e "9/(2+1)"`
---
### Imports
The syntax to use import is `import alias file_name.txt`
i.e
`import bills bills.txt` `import prob functions/probs.txt`

To use a variable from another file imported you must specify the `alias.` and then use variable name 
i.e
`free=bills.salary-2000` `print(foods.broccoli_kcal)`

i.e in a file
```
import tx imp/taxes.txt
import f finance.txt

print(tx.discount(f.salary))
```