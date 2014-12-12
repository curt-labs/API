GENERATE FILES:

PartNumbersNeeded - Lists old part numbers for which there is no new part number match

BaseVehiclesNeeded - Lists AAIA Base Vehicles for which there is no match in the BaseVehicles table 

SubmodelsNeeded - Lists AAIA Submodels for which there is no match in the Submodels table

NeedBaseVehicleInVcdbVehicleTable - Lists BaseVehicles that are needed in the vcdb_Vehicle table. That is, the BaseVehicle exists, but we need a Vehicle with VehicleID & BaseVehicleID, but with a SubmodelID and ConfigID of 0 in the vcdb_Vehicle table.

NeedSubmodelInVcdbVehicleTable - Lists Submodels that are needed in the vcdb_Vehicle table. That is, the Submodel exists, but we need a Vehicle with VehicleID,BaseModelID, & SubmodelID, but with a ConfigID of 0 in the vcdb_Vehicle table.

ConfigsNeeded - Giant list of all configs for all vehicles for which there is no match in the CurtDev database. 

ConfigsDiff - List of configs for all vehicles for which there is no match in the CurtDev database. 




METHODS:

Run - Starts program

CaptureCsv - Extracts data from the Csv file provided. Creates PartNumbersNeeded, BaseVehiclesNeeded, and SubmodelsNeeded files and writes to them. Returns a list of CsvData. Creates a map of CsvData by AAIABaseVehicleID

Audit - Begins the diff process. Passes 'map of CsvData by AAIABaseVehicleID' to AuditBaseVehicle. Passes returned Submodel map to AuditSubmodel method. Passes returned vehicle array to HandleVehicles method.

AuditBaseVehicles - For each groups of Vehicles (grouped by BaseVehicleID), looks at all parts. If all parts are the same for a BaseVehicle, it checks the VehcilePart table to see if the part is there. If not, it adds the part (commented out). If a basevehicle has different parts, it groups these vehicles by Submodel (submodelmap) and passes the map to AuditSubmodels.

AuditSubmodels - For each groups of Vehicles (grouped by Submodel), looks at all parts. If all parts are the same for a Submodel, it checks the VehcilePart table to see if the part is there. If not, it adds the part (commented out). If a Submodel has different parts, it adds these vehicles to a vehicle array and passes to HandleVehicles.

HandleVehicles - Takes an array of vehicles (as CsvData). Loops over vehicles. Loops over the ConfigAttributeTypes for which we have a corresponding entry in the ConfigAttributeType table. Assigns CurtDev ConfigAttributeTypeID' and ConfigAttributeID's to each configuration of each vehicle. Writes configs for which there is no attribute match to the ConfigsNeeded file. Creates a map of Vehicles by AAIAVehicleID. Calls diffVehicleConfigs.

diffVehicleConfigs - checks vehicles as grouped by VehicleID. If all vehicleConfigs are the same, the method checks & assigns the part to the Submodel (in vcdb_Vehcile and vcdb_VehiclePart). If configs vary for a VehicleID group, it checks for existing vehicles, inserts the vehicleConfig (as a vehicle in vcdb_Vehicle), and inserts the part.

