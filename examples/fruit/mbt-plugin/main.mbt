/// `void_func` - This demonstrates how you can create an export with
/// no inputs or outputs.
pub fn void_func() -> Unit {
  @host.log_debug_str("ENTER MoonBit plugin void_func")
  // TODO: fill out your implementation here
  @host.log_debug_str("LEAVE MoonBit plugin void_func")
}

/// `primitive_type_func` - This demonstrates how you can accept or return primtive types.
/// This function takes a utf8 string and returns a json encoded boolean
///
/// `input` - A string passed into plugin input
/// Returns A boolean encoded as json
pub fn primitive_type_func(input : String) -> Bool {
  @host.log_debug_str("ENTER MoonBit plugin primitive_type_func('\{input}')")
  // TODO: fill out your implementation here
  @host.log_debug_str("LEAVE MoonBit plugin primitive_type_func('\{input}')")
  false
}

/// `reference_type_func` - This demonstrates how you can accept or return references to schema types.
/// And it shows how you can define an enum to be used as a property or input/output.
pub fn reference_type_func(input : Fruit) -> ComplexObject {
  @host.log_debug_str("ENTER MoonBit plugin reference_type_func('\{input}')")
  // TODO: fill out your implementation here
  @host.log_debug_str("LEAVE MoonBit plugin reference_type_func('\{input}')")
  {
    ghost: Blinky,
    a_boolean: true,
    a_string: input.to_string(),
    an_int: 42,
    an_optional_date: Some("tbd"),
  }
}

fn main {

}
