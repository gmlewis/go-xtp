test "Address.to_json and .from_json work as expected on default object" {
  let default_object = Address::new()
  let got = default_object.to_json().stringify(escape_slash=false)
  let want =
    #|{"street":""}
  assert_eq!(got, want)
  //
  let got_parse : Address = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, default_object)
}

test "User.to_json and .from_json work as expected on default object" {
  let default_object = User::new()
  let got = default_object.to_json().stringify(escape_slash=false)
  let want =
    #|{}
  assert_eq!(got, want)
  //
  let got_parse : User = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, default_object)
}

test "User.to_json and .from_json work as expected on object only containing required fields" {
  let required_fields : User = {
    age: None,
    email: None,
    address: None,
  }
  let got = required_fields.to_json() |> @json.stringify(escape_slash=false)
  let want =
    #|{}
  assert_eq!(got, want)
  //
  let got_parse : User = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, required_fields)
}

test "User.to_json and .from_json work as expected on object with optional fields" {
  let optional_fields : User = {
    ..User::new(),
    age: Some(42),
    email: Some("email"),
    address: Some({street: ""}),
  }
  let got = optional_fields.to_json() |> @json.stringify(escape_slash=false)
  let want =
    #|{"age":42,"email":"email","address":{"street":""}}
  assert_eq!(got, want)
  //
  let got_parse : User = @json.from_json!(@json.parse!(want))
  assert_eq!(got_parse, optional_fields)
}
