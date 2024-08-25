pub fn host_eat_a_fruit(offset : Int64) -> Int64 = "extism:host/user" "eatAFruit"

/// `eat_a_fruit` - This is a host function. Right now host functions can only be the type (i64) -> i64.
/// We will support more in the future. Much of the same rules as exports apply.
pub fn eat_a_fruit(input : Fruit) -> Bool!String {
  let json = @jsonutil.to_json(input)
  let mem = @host.Memory::allocate_json_value(json)
  let ptr = host_eat_a_fruit(mem.offset)
  let buf = @host.find_memory(ptr).to_string()
  let out = @json.parse(buf)!!
  match out {
    Ok(jv) =>
      match jv.as_bool() {
        Some(v) => v
        None => {
          raise "unable to parse \{buf}"
        }
      }
    Err(e) => {
      raise "unable to parse \{buf}: \{e}"
    }
  }
}
