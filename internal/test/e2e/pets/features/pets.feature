Feature:
  In order to use pets api
  As an user of the Api peya
  I need to able to manage pets

  Scenario: get all pets registered
    Given the are registered pets
      | id | name  | vaccines             | age_months |
      | 1  | Rino  | distemper,parvovirus | 36         |
      | 2  | Braco | rabies               | 3          |
      | 3  | Duke  |                      | 12         |
    When I send "GET" request to "/pets"
    Then the response status code should be 200
    And the response should match json:
      """
        [
          {
            "id":1,
            "name":"Rino",
            "vaccines":[
               "distemper",
               "parvovirus"
            ],
            "age": "3 years"
          },
          {
            "id":2,
            "name":"Braco",
            "vaccines":[
               "rabies"
            ],
            "age": "3 months"
          },
             {
            "id":3,
            "name":"Duke",
            "vaccines":[ ],
            "age": "1 year"
          }
        ]
      """