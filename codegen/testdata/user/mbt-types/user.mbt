/// `Address` represents a users address.
pub struct Address {
  /// Street address
  street : String
} derive(Show, Eq)

/// `Address::new` returns a new struct with default values.
pub fn Address::new() -> Address {
  {
    street: "",
  }
}

pub impl @jsonutil.ToJson for Address with to_json(self) {
  let fields : Array[(String, @jsonutil.ToJson)] = [
    ("street", self.street),
  ]
  @jsonutil.from_entries(fields)
}

/// `Address::from_json` transforms a `Json` to a value.
pub fn Address::from_json(value : Json) -> Address? {
  let value = value.as_object()?
  let street = value.get("street")?.as_string()
  match (
    street,
  ) {
      Some(street)
    => Some({
      street,
    })
    _ => None
  }
}

/// `Address::parse` parses a JSON string and returns the value.
pub fn Address::parse(s : String) -> Address!String {
  match @json.parse(s)!! {
    Ok(jv) =>
      match Address::from_json(jv) {
        Some(value) => value
        None => {
          raise "unable to parse Address \(s)"
        }
      }
    Err(e) => {
      raise "unable to parse Address \(s): \(e)"
    }
  }
}

/// `get_schema` returns an `XTPSchema` for the `Address`.
pub fn get_schema(self : Address) -> XTPSchema {
  {
    "street": "string",
  }
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
  {
    age: None,
    email: None,
    address: None,
  }
}

pub impl @jsonutil.ToJson for User with to_json(self) {
  let fields : Array[(String, @jsonutil.ToJson)] = [
  ]
  match self.age {
    Some(value) => fields.append([("age", value)])
    None => ()
  }
  match self.email {
    Some(value) => fields.append([("email", value)])
    None => ()
  }
  match self.address {
    Some(value) => fields.append([("address", value)])
    None => ()
  }
  @jsonutil.from_entries(fields)
}

/// `User::from_json` transforms a `Json` to a value.
pub fn User::from_json(value : Json) -> User? {
  let value = value.as_object()?
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
  match (
    age,
    email,
    address,
  ) {
    (
      age,
      email,
      address,
    ) => Some({
      age,
      email,
      address,
    })
    _ => None
  }
}

/// `User::parse` parses a JSON string and returns the value.
pub fn User::parse(s : String) -> User!String {
  match @json.parse(s)!! {
    Ok(jv) =>
      match User::from_json(jv) {
        Some(value) => value
        None => {
          raise "unable to parse User \(s)"
        }
      }
    Err(e) => {
      raise "unable to parse User \(s): \(e)"
    }
  }
}

/// `get_schema` returns an `XTPSchema` for the `User`.
pub fn get_schema(self : User) -> XTPSchema {
  {
    "age": "?integer",
    "email": "?string",
    "address": "?Address",
  }
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
