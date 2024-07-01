/// Exported: voidFunc
pub fn exported_void_func() -> Int {
  void_func()
  return 0 // success
}

/// Exported: primitiveTypeFunc
pub fn exported_primitive_type_func() -> Int {
  let input = @host.input_string()
  let output = primitive_type_func(input) |> @jsonutil.to_json()
  @host.output_json_value(output)
  return 0 // success
}

/// Exported: referenceTypeFunc
pub fn exported_reference_type_func() -> Int {
  Fruit::parse(@host.input_string())!!.unwrap()
  |> reference_type_func()
  |> @jsonutil.to_json()
  |> @host.output_json_value()
  return 0 // success
}
