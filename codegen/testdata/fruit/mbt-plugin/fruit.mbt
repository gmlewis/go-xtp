/// `Fruit` represents a set of available fruits you can consume.
pub enum Fruit {
  Apple
  Orange
  Banana
  Strawberry
} derive(Eq)

// Why is `Fruit.to_string` necessary when the `Show` trait is implemented below?
pub fn to_string(self : Fruit) -> String {
  let buf = Buffer::new()
  Show::output(self, buf)
  buf.to_string()
}

/// `Fruit.output` implements the Show trait.
pub impl Show for Fruit with output(self, logger) {
  match self {
    Apple => logger.write_string("apple")
    Orange => logger.write_string("orange")
    Banana => logger.write_string("banana")
    Strawberry => logger.write_string("strawberry")
  }
}

pub fn to_json(self : Fruit) -> Json {
  match self {
    Apple => "apple".to_json()
    Orange => "orange".to_json()
    Banana => "banana".to_json()
    Strawberry => "strawberry".to_json()
  }
}

/// `Fruit::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for Fruit with from_json(json, path) {
  match json {
    String("apple") => Apple
    String("orange") => Orange
    String("banana") => Banana
    String("strawberry") => Strawberry
    s =>
      raise @json.JsonDecodeError(
        (path, "Fruit::from_json: expected a Fruit, got \{s}"),
      )
  }
}

/// `GhostGang` represents a set of all the enemies of pac-man.
pub enum GhostGang {
  Blinky
  Pinky
  Inky
  Clyde
} derive(Eq)

// Why is `GhostGang.to_string` necessary when the `Show` trait is implemented below?
pub fn to_string(self : GhostGang) -> String {
  let buf = Buffer::new()
  Show::output(self, buf)
  buf.to_string()
}

/// `GhostGang.output` implements the Show trait.
pub impl Show for GhostGang with output(self, logger) {
  match self {
    Blinky => logger.write_string("blinky")
    Pinky => logger.write_string("pinky")
    Inky => logger.write_string("inky")
    Clyde => logger.write_string("clyde")
  }
}

pub fn to_json(self : GhostGang) -> Json {
  match self {
    Blinky => "blinky".to_json()
    Pinky => "pinky".to_json()
    Inky => "inky".to_json()
    Clyde => "clyde".to_json()
  }
}

/// `GhostGang::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for GhostGang with from_json(json, path) {
  match json {
    String("blinky") => Blinky
    String("pinky") => Pinky
    String("inky") => Inky
    String("clyde") => Clyde
    s =>
      raise @json.JsonDecodeError(
        (path, "GhostGang::from_json: expected a GhostGang, got \{s}"),
      )
  }
}

/// `ComplexObject` represents a complex json object.
pub struct ComplexObject {
  /// I can override the description for the property here
  ghost : GhostGang
  /// A boolean prop
  a_boolean : Bool
  /// An string prop
  a_string : String
  /// An int prop
  an_int : Int
  /// A datetime object, we will automatically serialize and deserialize
  /// this for you.
  an_optional_date : String?
} derive(Show, Eq)

/// `ComplexObject::new` returns a new struct with default values.
pub fn ComplexObject::new() -> ComplexObject {
  {
    ghost: Blinky,
    a_boolean: false,
    a_string: "",
    an_int: 0,
    an_optional_date: None,
  }
}

pub fn to_json(self : ComplexObject) -> Json {
  let json : Map[String, Json] = {  }
  json["ghost"] = self.ghost.to_json()
  json["aBoolean"] = self.a_boolean.to_json()
  json["aString"] = self.a_string.to_json()
  json["anInt"] = self.an_int.to_json()
  match self.an_optional_date {
    Some(an_optional_date) =>
      json["anOptionalDate"] = an_optional_date.to_json()
    _ => ()
  }
  json.to_json()
}

/// `ComplexObject::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for ComplexObject with from_json(json, path) {
  let json = match json.as_object() {
    Some(json) => json
    e =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json: expected object, got \{e}"),
      )
  }
  let ghost : GhostGang = match json.get("ghost") {
    Some(ghost) => @json.from_json!(ghost)
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json:ghost: expected GhostGang"),
      )
  }
  let a_boolean : Bool = match json.get("aBoolean") {
    Some(True) => true
    Some(False) => false
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json:a_boolean: expected Bool"),
      )
  }
  let a_string : String = match json.get("aString") {
    Some(String(a_string)) => a_string
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json:a_string: expected String"),
      )
  }
  let an_int : Int = match json.get("anInt") {
    Some(Number(an_int)) => an_int.to_int()
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json:an_int: expected Int"),
      )
  }
  let an_optional_date : String? = match json.get("anOptionalDate") {
    Some(String(an_optional_date)) => Some(an_optional_date)
    Some(Null) | None => None
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json:an_optional_date: expected String? or Null"),
      )
  }
  {
    ghost,
    a_boolean,
    a_string,
    an_int,
    an_optional_date,
  }
}

/// `ComplexObject::get_schema` returns an `XTPSchema` for the `ComplexObject`.
pub fn ComplexObject::get_schema() -> XTPSchema {
  {
    "ghost": "GhostGang",
    "aBoolean": "boolean",
    "aString": "string",
    "anInt": "integer",
    "anOptionalDate": "?Date",
  }
}

/// `XTPSchema` describes the values and types of an XTP object
/// in a language-agnostic format.
type XTPSchema Map[String, String]
