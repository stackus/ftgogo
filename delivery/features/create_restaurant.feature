@command @restaurant
Feature: Create Restaurants

  Scenario: Can create new restaurants
    When I create a restaurant named "Best Foods" with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to succeed

  Scenario: Creating duplicate restaurants returns an error
    Given I create a restaurant named "Best Foods" with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    When I create another restaurant named "Best Foods" with address
      | Street1 | 123 Address St. |
      | City    | BigTown         |
      | State   | Colorado        |
      | Zip     | 80120           |
    Then I expect the command to fail
    And the returned error message is "restaurant already exists"
