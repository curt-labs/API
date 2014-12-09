

Polk CSV provides vcdb VehicleID, BaseID, SubID and several Aces Configs


Diff/Insert workflow


Check for CurtVehicleID by AcesVehicleID + AcesConfigs //We have this exact vehicle

	If, match, add part to vcdb_VehicleParts
	return

	
Check for CurtVehicle by AcesVehicleID //Config needs to be added

	Add Curt Configs
	Add part to vcdb_VehicleParts
	return 


Check for BaseVehicle + Submodel By AcesBaseVehicle + Submodel //we have these basics
	Add Curt Configs
	Add part to vcdb_VehicleParts
	return 

Add Base + Submodel + Configs
	Add part to vcdb_VehicleParts
	return 



//REVISED

1) Parse Polk CSV in structs, CsvDatum

2) Loop through []CsvDatum and put them in maps of:
		[basevehicleID][]CsvDatum
		[submodelID][]CsvDatum

3) Loop through [basevehicleID][]CsvDatum
	For each basevehicleID:
	a) If all part numbers are the same:
		**TODO - convert Aries Part Numbers
		**TODO - check to see if vcdb_VehiclePart exists. If not:
		Insert part and basevehicle into vcdb_VehiclePart table (vcdb_Vehicle where submodel == 0)
		if finding this vcdb_Vehicle fails, insert it, then insert vcdb_VehiclePart

	b) If there are different part numbers within the basevehicle map, loop through [submodelID][]CsvDatum
		For each submodelID:
		i) If all part numbers are the same:
			**TODO - convert Aries Part Numbers
			**TODO - check to see if vcdb_VehiclePart exists. If not:
			Insert part and submodel into vcdb_VehiclePart table (vcdb_Vehicle where ConfigID == 0)
			if finding this vcdb_Vehicle fails, insert it, then insert vcdb_VehiclePart

		ii) If there are different part numbers within the submodel map, ...TODO 



		**TODO - work through map by config