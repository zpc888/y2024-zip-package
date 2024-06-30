#![allow(unused)]
#[repr(transparent)]
#[derive(Debug)]
struct CheckNumber(u32);
#[repr(transparent)]
#[derive(Debug)]
struct CardNumber(String);

enum CardType {
    Visa,
    MasterCard,
    Amex,
}

#[derive(Debug)]
struct CreditCard(CheckNumber, CardNumber);

#[derive(Debug)]
enum PaymentMethod {
    Cash,
    Check(CheckNumber),
    Card(CreditCard),
}

#[derive(Debug)]
#[repr(transparent)]
struct PaymentAmountInCent(i32);

#[derive(Debug)]
enum Currency {
    Eur,
    Usd,
    Cad,
    Gbp,
}

#[derive(Debug)]
struct Payment {
    amount_in_cent: PaymentAmountInCent,
    currency: Currency,
    method: PaymentMethod,
}

trait PrintDetails {
    fn payment_info(&self) -> String;
}

impl PrintDetails for Payment {
    fn payment_info(&self) -> String {
        let method = match &self.method {
            PaymentMethod::Cash => String::from("cash"),
            PaymentMethod::Check(check_number) => {
                format!("a check with number {:?}", check_number.0)
            }
            PaymentMethod::Card(credit_card) => format!(
                "a credit card {} with check number {}",
                credit_card.1 .0, credit_card.0 .0
            ),
        };

        return format!(
            "An amount of {} in cents, was paid in {:?} using {}",
            &self.amount_in_cent.0, &self.currency, method
        );
    }
}

fn main2() {
    let cc_payment = Payment {
        amount_in_cent: PaymentAmountInCent(1000),
        currency: Currency::Usd,
        method: PaymentMethod::Card(CreditCard(
            CheckNumber(88866),
            CardNumber("1234 5678 9012 3456".to_string()),
        )),
    };
    println!(
        "credit card payment details:\n {}",
        cc_payment.payment_info()
    );

    let check_payment = Payment {
        amount_in_cent: PaymentAmountInCent(800),
        currency: Currency::Cad,
        method: PaymentMethod::Check(CheckNumber(123456)),
    };
    println!("check payment details:\n {}", check_payment.payment_info());
}
