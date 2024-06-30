/// Exported: voidFunc
pub fn voidFunc() -> Int {
  void_func()
  return 0 // success
}

/// Exported: primitiveTypeFunc
pub fn primitiveTypeFunc() -> Int {
  let input = @host.input_string()
  let output = primitive_type_func(input) |> @jsonutil.to_json()
  @host.output_json_value(output)
  return 0 // success
}

/// Exported: referenceTypeFunc
pub fn referenceTypeFunc() -> Int {
  Fruit::parse(@host.input_string())!!.unwrap()
  |> reference_type_func()
  |> @jsonutil.to_json()
  |> @host.output_json_value()
  return 0 // success
}
