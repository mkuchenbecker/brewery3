docker run -d mkuchenbecker/brewery3:thermometer-latest --env THERM_PORT="9110" --env THERM_ADDR="28-0315712c08ff"

	strPort := os.Getenv("THERM_PORT")
	address := os.Getenv("THERM_ADDR")

            name: THERM_PORT="9110"
          - name: THERM_ADDR= "28-0315712c08ff"


          - name: THERM_PORT="9111"
          - name: THERM_ADDR= "28-0315715039ff"


          - name: THERM_PORT="9112"
          - name: THERM_ADDR= "28-031571188aff"