#![allow(unused_imports)]

use std::collections::HashMap;
use std::thread::sleep;
use std::time::Duration;

mod algebraic_type_playground;

#[tokio::main]
async fn main() {
    println!("Rust combines many languages features in a unique way");
    println!("\t1. C++ RAII for ownership/borrow and Zero-Cost Abstraction");
    println!("\t2. ADT from Haskell");
    println!("\t3. Polymorphism via traits");
    println!("\t4. Async/await from JavaScript");
    println!("\t5. Meta Programming - Macro");
    println!("\t... ...");
    println!("==============Zeros-Cost Abstraction================================================");
    iter_abstraction();
    println!("==============ADT from Haskell======================================================");
    algebraic_data_type();
    println!("==============Polymorphism via trait================================================");
    polymorphism();
    println!("==============Async/await from JavaScript===========================================");
    async_await().await;
    println!("==============Meta Programming - Macro===========================================");
    meta_programming_macro();

    // sleep(Duration::from_secs(1));
}

type EmployeeName = String;

enum Employee {
    Manager {
        name: EmployeeName,
        subordinates: Vec<Box<Employee>>,
    },
    Worker {
        name: EmployeeName,
        manager: String,
    },
}

fn print_employee_type(employee: &Employee) {
    match employee {
        Employee::Manager { name, subordinates } => {
            println!("Manager: {} with {} subordinates", name, subordinates.len());
        }
        Employee::Worker { name, manager } => {
            println!("Worker: {} managed by {}", name, manager);
        }
    }
}

fn algebraic_data_type() {
    // SUM type  vs PRODUCT type
    let bob = Employee::Worker {
        name: "Bob".to_string(),
        manager: "Alice".to_string(),
    };
    let charles = Employee::Worker {
        name: "Charles".to_string(),
        manager: "Alice".to_string(),
    };
    print_employee_type(&bob);
    print_employee_type(&charles);
    let alice = Employee::Manager {
        name: "Alice".to_string(),
        subordinates: vec![Box::new(bob), Box::new(charles)],
    };
    let employees = vec![&alice];
    for employee in employees {
        print_employee_type(employee)
    }
}

fn iter_abstraction() {
    let numbers = vec![5, 4, 9, 3, 2, 1, 6, 7, 8];

    // using iterator abstraction to find max - ZERO COST ABSTRACTION
    let max_elem = numbers.iter().max();

    match max_elem {
        Some(&max) => println!("Max number is: {}", max),
        None => println!("No max number found"),
    }
}

trait Shape {
    fn area(&self) -> f64;
}

struct Rectangle {
    width: f64,
    height: f64,
}

struct Circle {
    radius: f64,
}

impl Shape for Rectangle {
    fn area(&self) -> f64 {
        self.width * self.height
    }
}

impl Shape for Circle {
    fn area(&self) -> f64 {
        std::f64::consts::PI * self.radius * self.radius
    }
}

fn print_area<T: Shape>(shape: &T) {
    println!("Area: {}", shape.area());
}

fn print_area_dyn(shape: &dyn Shape) {
    println!("Area: {}", shape.area());
}

fn polymorphism() {
    let rectangle = Rectangle {
        width: 3.0,
        height: 4.0,
    };
    let circle = Circle { radius: 5.0 };

    print_area(&rectangle);
    print_area(&circle);
    println!(
        "rectangle width: {} * height: {} = area: {}",
        rectangle.width,
        rectangle.height,
        rectangle.area()
    );
    println!("circle radius: {} = area: {}", circle.radius, circle.area());
/*
    let shapes: Vec<Box<dyn Shape>> = vec![Box::new(rectangle), Box::new(circle)];
    for shape in shapes {
        println!("Area: {}", shape.area());
    }
 */

    let shapes: Vec<Box<dyn Shape>> = vec![Box::new(rectangle), Box::new(circle)];
    for shape in shapes {
        print_area_dyn(&*shape);
    }
}

async fn async_await() {
    println!("making httpbin call");
    let delay = 5;
    // let response = reqwest::get(&format!("http://httpbin.org/delay/{delay}"))
    //     .await.unwrap().text().await;
    let client = reqwest::Client::new();
    let response = client.get(&format!("http://httpbin.org/delay/{delay}")).send()
        .await.unwrap();
    println!("Response JSON: {:?}", response);
}

fn meta_programming_macro() {
    macro_rules! map {
        ($key:ty, $val:ty) => {
          {
              let map: HashMap<$key, $val> = HashMap::new();
              map
          }
        };
        ($($key:expr => $value:expr),*) => {
            {
                let mut map = HashMap::new();
                $(
                    map.insert($key, $value);
                )*
                map
            }
        };
    }

    let scores_without_macro: HashMap<String, i32> = HashMap::new();
    let scores_with_macro = map!(String, i32);

    let mut mut_without_macro = HashMap::new();
    mut_without_macro.insert("Blue".to_string(), 3);
    mut_without_macro.insert("Red".to_string(), 5);
    mut_without_macro.insert("Green".to_owned(), 1);        // what is difference between to_owned vs to_string

    let mut_with_macro = map!{
        "Blue".to_string() => 3,
        "Red".to_string() => 5,
        "Green".to_string() => 1
    };
    println!("scores_without_macro: {:?}", scores_without_macro);
    println!("scores_with_macro: {:?}", scores_with_macro);
    println!("mut_scores_without_macro: {:?}", mut_without_macro);
    println!("mut_scores_with_macro: {:?}", mut_with_macro);
}
