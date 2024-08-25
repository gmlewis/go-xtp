/// `Fruit` represents a set of available fruits you can consume.
pub enum Fruit {
  Apple
  Orange
  Banana
  Strawberry
} derive(Show, Eq)

/// `Fruit.output` implements the Show trait.
pub fn output(self : Fruit) -> String {
  match self {
    Apple => "apple"
    Orange => "orange"
    Banana => "banana"
    Strawberry => "strawberry"
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
        (path, "Fruit::from_json expected a Fruit, got \{s}"),
      )
  }
}

/// `GhostGang` represents a set of all the enemies of pac-man.
pub enum GhostGang {
  Blinky
  Pinky
  Inky
  Clyde
} derive(Show, Eq)

/// `GhostGang.output` implements the Show trait.
pub fn output(self : GhostGang) -> String {
  match self {
    Blinky => "blinky"
    Pinky => "pinky"
    Inky => "inky"
    Clyde => "clyde"
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
        (path, "GhostGang::from_json expected a GhostGang, got \{s}"),
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
  json["a_boolean"] = self.a_boolean.to_json()
  json["a_string"] = self.a_string.to_json()
  json["an_int"] = self.an_int.to_json()
  match self.an_optional_date {
    Some(an_optional_date) =>
      json["an_optional_date"] = an_optional_date.to_json()
    _ => ()
  }
  json.to_json()
}

/// `ComplexObject::from_json` transforms a `Json` to a value.
pub impl @json.FromJson for ComplexObject with from_json(json, path) {
  let json = match json.as_object() {
    Some(json) => json
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplexObject::from_json: expected object"),
      )
  }
  let ghost : GhostGang = match json.get("ghost") {
    Some(ghost) => @json.from_json!(ghost)
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplextObject::from_json:ghost: expected GhostGang"),
      )
  }
  let a_boolean : Bool = match json.get("a_boolean") {
    Some(True) => true
    Some(False) => false
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplextObject::from_json:a_boolean expected Bool"),
      )
  }
  let a_string : String = match json.get("a_string") {
    Some(String(a_string)) => a_string
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplextObject::from_json:a_string: expected String"),
      )
  }
  let an_int : Int = match json.get("an_int") {
    Some(Number(an_int)) => an_int.to_int()
    _ =>
      raise @json.JsonDecodeError(
        (path, "ComplextObject::from_json:an_int expected Int"),
      )
  }
  let an_optional_date : String? = match json.get("an_optional_date") {
    Some(String(an_optional_date)) => Some(an_optional_date)
    Some(Null) | None => None
    _ =>
      raise @json.JsonDecodeError(
        (
          path, "ComplextObject::from_json:an_optional_date expected String or Null",
        ),
      )
  }
  { ghost, a_boolean, a_string, an_int, an_optional_date }
}

/// `get_schema` returns an `XTPSchema` for the `ComplexObject`.
pub fn get_schema(self : ComplexObject) -> XTPSchema {
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
