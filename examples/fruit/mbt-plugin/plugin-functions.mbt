/// Exported: voidFunc
pub fn exported_void_func() -> Unit {
	void_func()
}

/// Exported: primitiveTypeFunc
pub fn exported_primitive_type_func() -> Int {
	let input = @host.InputString()
	let output = primitive_type_func(input)

	let buf, err = json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}

/// Exported: referenceTypeFunc
pub fn exported_reference_type_func() -> Int {
	let input = pdk.InputString()
	let output = reference_type_func(Fruit(input))

	let buf, err = json.Marshal(output)
	if err != nil {
		pdk.Log(pdk.LogError, err.Error())
		return 1 // failure
	}

	pdk.OutputString(string(buf))
	return 0 // success
}

