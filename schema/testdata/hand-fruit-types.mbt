/// `Fruit` represents a set of available fruits you can consume.
pub enum Fruit {
  Apple
  Orange
  Banana
  Strawberry
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
}

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
