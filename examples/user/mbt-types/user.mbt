/// `Address` represents a users address.
pub struct Address {
  /// Street address
  street : String
} derive(Show, Eq)

/// `Address::new` returns a new struct with default values.
pub fn Address::new() -> Address {
  { street: "" }
}

pub fn to_json(self : Address) -> Json {
  { "street": self.street.to_json() }
}

/// `Address::from_json` transforms a `Json` to a value.
pub fn Address::from_json(value : Json) -> Address? {
  let value = match value.as_object() {
    Some(v) => v
    _ => return None
  }
  let street = match value.get("street") {
    Some(v) => v
    _ => return None
  }
  let street = street.as_string()
  match street {
    Some(street) => Some({ street, })
    _ => None
  }
}

pub type! AddressError String derive(Show)

/// `Address::parse` parses a JSON string and returns the value.
pub fn Address::parse(s : String) -> Address!AddressError {
  match @json.parse?(s) {
    Ok(jv) =>
      match Address::from_json(jv) {
        Some(value) => value
        None => raise AddressError("unable to parse Address \{s}")
      }
    Err(e) => raise AddressError("unable to parse Address \{s}: \{e}")
  }
}

/// `get_schema` returns an `XTPSchema` for the `Address`.
pub fn get_schema(self : Address) -> XTPSchema {
  { "street": "string" }
}

/// `User` represents a user object in our system..
pub struct User {
  /// The user's age, naturally
  age : Int?
  /// The user's email, of course
  email : String?
  address : Address?
} derive(Show, Eq)

/// `User::new` returns a new struct with default values.
pub fn User::new() -> User {
  { age: None, email: None, address: None }
}

pub fn to_json(self : User) -> Json {
  let object = {  }
  match self.age {
    Some(value) => object.set("age", value.to_json())
    None => ()
  }
  match self.email {
    Some(value) => object.set("email", value.to_json())
    None => ()
  }
  match self.address {
    Some(value) => object.set("address", value.to_json())
    None => ()
  }
  object.to_json()
}

/// `User::from_json` transforms a `Json` to a value.
pub fn User::from_json(value : Json) -> User? {
  let value = match value.as_object() {
    Some(v) => v
    _ => return None
  }
  let age = match value.get("age") {
    Some(jv) => json_as_integer(jv)
    None => None
  }
  let email = match value.get("email") {
    Some(jv) => jv.as_string()
    None => None
  }
  let address = match value.get("address") {
    Some(jv) => Address::from_json(jv)
    None => None
  }
  Some({ age, email, address })
}

pub type! UserError String derive(Show)

/// `User::parse` parses a JSON string and returns the value.
pub fn User::parse(s : String) -> User!UserError {
  match @json.parse?(s) {
    Ok(jv) =>
      match User::from_json(jv) {
        Some(value) => value
        None => raise UserError("unable to parse User \{s}")
      }
    Err(e) => raise UserError("unable to parse User \{s}: \{e}")
  }
}

/// `get_schema` returns an `XTPSchema` for the `User`.
pub fn get_schema(self : User) -> XTPSchema {
  { "age": "?integer", "email": "?string", "address": "?Address" }
}

/// `XTPSchema` describes the values and types of an XTP object
/// in a language-agnostic format.
type XTPSchema Map[String, String]

// https://github.com/moonbitlang/core/issues/651
fn json_as_integer(value : Json) -> Int? {
  match value.as_number() {
    Some(n) => Some(n.to_int())
    None => None
  }
}
