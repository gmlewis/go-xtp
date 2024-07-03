/// `Address` represents a users address.
pub struct Address {
  /// Street address
  street : String
} derive(Debug, Eq)

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

pub fn Address::from_json(value : @json.JsonValue) -> Address? {
  match value {
    @json.JsonValue::Object({
      "street": Some(@json.JsonValue::String(s)),
    }) => Some({
      street: s,
    })
    _ => None
  }
}

pub fn Address::parse(s : String) -> Address!String {
  match @json.parse(s) {
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
} derive(Debug, Eq)

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

pub fn User::from_json(value : @json.JsonValue) -> User? {
  match value.as_object() {
    Some(m) =>
      Some(
        {
          age: match m.get_or_default("age", @json.JsonValue::Null) {
            @json.JsonValue::Number(n) => Some(n.to_int())
            _ => None
          },
          email: match m.get_or_default("email", @json.JsonValue::Null) {
            @json.JsonValue::String(s) => Some(s.to_string())
            _ => None
          },
          address: m.get_or_default("address", @json.JsonValue::Null)
            |> Address::from_json(),
        },
      )
    None => None
  }
}

pub fn User::parse(s : String) -> User!String {
  match @json.parse(s) {
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

// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema Map[String, String]
