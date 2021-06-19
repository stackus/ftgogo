@query @address
Feature: Get Consumer Address

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Can get a consumers address
    Given I add an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "home",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    When I request an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "home"
    }
    """
    Then I expect the request to succeed
    And the returned address matches:
    """
    {
      "Street1": "123 Address St.",
      "Street2": "",
      "City": "HomeTown",
      "State": "HomeState",
      "Zip": "12345"
    }
    """

  Scenario: Getting an address that doesn't exist returns an error
    Given I add an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "home",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    When I request an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "other"
    }
    """
    Then I expect the request to fail
    And the returned error message is:
    """
    an address with that identifier does not exist
    """

  Scenario: Getting an address for a consumer that doesn't exist returns an error
    Given I add an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "home",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    When I request an address with:
    """
    {
      "ConsumerID": "b456",
      "AddressID": "home"
    }
    """
    Then I expect the request to fail
    And the returned error message is:
    """
    consumer not found
    """
