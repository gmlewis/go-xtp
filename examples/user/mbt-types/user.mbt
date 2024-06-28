/// `Address` represents a users address.
pub struct Address {
  /// Street address
  street : String
}

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

/// `User` represents a user object in our system..
pub struct User {
  /// The user's age, naturally
  age : Int?
  /// The user's email, of course
  email : String?
  address : Address?
}

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
