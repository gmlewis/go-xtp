/// `void_func` - This demonstrates how you can create an export with
/// no inputs or outputs.
pub fn void_func() -> Unit {
  // TODO: fill out your implementation here
}

/// `primitive_type_func` - This demonstrates how you can accept or return primtive types.
/// This function takes a utf8 string and returns a json encoded boolean
//
/// `input` - A string passed into plugin input
/// Returns A boolean encoded as json
pub fn primitive_type_func(input : String) -> Bool {
  // TODO: fill out your implementation here
  false
}

/// `reference_type_func` - This demonstrates how you can accept or return references to schema types.
/// And it shows how you can define an enum to be used as a property or input/output.
pub fn reference_type_func(input : Fruit) -> ComplexObject {
  { ..ComplexObject::new(), a_boolean: eat_a_fruit(input)!!.unwrap() }
}

fn main {

}
