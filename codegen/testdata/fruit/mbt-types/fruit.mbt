/// `Fruit` represents a set of available fruits you can consume.
pub enum Fruit {
  Apple
  Orange
  Banana
  Strawberry
} derive(Debug, Eq)

pub fn Fruit::parse(s : String) -> Fruit!String {
  match s {
    "apple" => Apple
    "orange" => Orange
    "banana" => Banana
    "strawberry" => Strawberry
    _ => {
      raise "not a Fruit: \(s)"
    }
  }
}

pub impl @jsonutil.ToJson for Fruit with to_json(self) {
  match self {
    Apple => @jsonutil.to_json("apple")
    Orange => @jsonutil.to_json("orange")
    Banana => @jsonutil.to_json("banana")
    Strawberry => @jsonutil.to_json("strawberry")
  }
}

/// `GhostGang` represents a set of all the enemies of pac-man.
pub enum GhostGang {
  Blinky
  Pinky
  Inky
  Clyde
} derive(Debug, Eq)

pub fn GhostGang::parse(s : String) -> GhostGang!String {
  match s {
    "blinky" => Blinky
    "pinky" => Pinky
    "inky" => Inky
    "clyde" => Clyde
    _ => {
      raise "not a GhostGang: \(s)"
    }
  }
}

pub impl @jsonutil.ToJson for GhostGang with to_json(self) {
  match self {
    Blinky => @jsonutil.to_json("blinky")
    Pinky => @jsonutil.to_json("pinky")
    Inky => @jsonutil.to_json("inky")
    Clyde => @jsonutil.to_json("clyde")
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
} derive(Debug, Eq)

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

pub impl @jsonutil.ToJson for ComplexObject with to_json(self) {
  let fields : Array[(String, @jsonutil.ToJson)] = [
    ("ghost", self.ghost),
    ("aBoolean", self.a_boolean),
    ("aString", self.a_string),
    ("anInt", self.an_int),
  ]
  match self.an_optional_date {
    Some(value) => fields.append([("anOptionalDate", value)])
    None => ()
  }
  @jsonutil.from_entries(fields)
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

// XTPSchema describes the values and types of an XTP object
// in a language-agnostic format.
type XTPSchema Map[String, String]
