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

pub fn to_json(self : Address) -> Json {
  let json : Map[String, Json] = {  }
  json["street"] = self.street.to_json()
  json.to_json()
}

/// `Address::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for Address with from_json(json, path) {
  let json = match json.as_object() {
    Some(json) => json
    e =>
      raise @json.JsonDecodeError(
        (path, "Address::from_json: expected object, got \{e}"),
      )
  }
  let street : String = match json.get("street") {
    Some(String(street)) => street
    _ =>
      raise @json.JsonDecodeError(
        (path, "Address::from_json:street: expected String"),
      )
  }
  {
    street,
  }
}

/// `Address::get_schema` returns an `XTPSchema` for the `Address`.
pub fn Address::get_schema() -> XTPSchema {
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

pub fn to_json(self : User) -> Json {
  let json : Map[String, Json] = {  }
  match self.age {
    Some(age) =>
      json["age"] = age.to_json()
    _ => ()
  }
  match self.email {
    Some(email) =>
      json["email"] = email.to_json()
    _ => ()
  }
  match self.address {
    Some(address) =>
      json["address"] = address.to_json()
    _ => ()
  }
  json.to_json()
}

/// `User::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for User with from_json(json, path) {
  let json = match json.as_object() {
    Some(json) => json
    e =>
      raise @json.JsonDecodeError(
        (path, "User::from_json: expected object, got \{e}"),
      )
  }
  let age : Int? = match json.get("age") {
    Some(Number(age)) => Some(age.to_int())
    Some(Null) | None => None
    _ =>
      raise @json.JsonDecodeError(
        (path, "User::from_json:age: expected Int? or Null"),
      )
  }
  let email : String? = match json.get("email") {
    Some(String(email)) => Some(email)
    Some(Null) | None => None
    _ =>
      raise @json.JsonDecodeError(
        (path, "User::from_json:email: expected String? or Null"),
      )
  }
  let address : Address? = match json.get("address") {
    Some(Object(address)) => Some(@json.from_json!(address.to_json()))
    Some(Null) | None => None
    _ =>
      raise @json.JsonDecodeError(
        (path, "User::from_json:address: expected Address? or Null"),
      )
  }
  {
    age,
    email,
    address,
  }
}

/// `User::get_schema` returns an `XTPSchema` for the `User`.
pub fn User::get_schema() -> XTPSchema {
  {
    "age": "?integer",
    "email": "?string",
    "address": "?Address",
  }
}

/// `XTPSchema` describes the values and types of an XTP object
/// in a language-agnostic format.
type XTPSchema Map[String, String]
