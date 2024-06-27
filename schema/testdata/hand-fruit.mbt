/// Fruit represents a set of available fruits you can consume.
pub enum Fruit {
  Apple
  Orange
  Banana
  Strawberry
}

pub impl @jsonutil.ToJson for Fruit with to_json(self) {
  match self {
    Apple => @json.JsonValue::String("apple")
    Orange => @json.JsonValue::String("orange")
    Banana => @json.JsonValue::String("banana")
    Strawberry => @json.JsonValue::String("strawberry")
  }
}

/// GhostGang represents a set of all the enemies of pac-man.
pub enum GhostGang {
  Blinky
  Pinky
  Inky
  Clyde
}

pub impl @jsonutil.ToJson for GhostGang with to_json(self) {
  match self {
    Blinky => @json.JsonValue::String("blinky")
    Pinky => @json.JsonValue::String("pinky")
    Inky => @json.JsonValue::String("inky")
    Clyde => @json.JsonValue::String("clyde")
  }
}

/// ComplexObject represents a complex json object.
pub struct ComplexObject {
  // I can override the description for the property here
  ghost : GhostGang
  // A boolean prop
  aBoolean : Bool
  // An string prop
  aString : String
  // An int prop
  anInt : Int
  // A datetime object, we will automatically serialize and deserialize
  // this for you.
  anOptionalDate : String?
}

pub impl @jsonutil.ToJson for ComplexObject with to_json(self) {
  let fields : Array[(String, @jsonutil.ToJson)] = [
    ("ghost", self.ghost),
    ("aBoolean", self.aBoolean),
    ("aString", self.aString),
    ("anInt", self.anInt),
  ]
  match self.anOptionalDate {
    Some(value) => fields.append(("anOptionalDate", value))
  }
  @jsonutil.from_entries(fields)
}
