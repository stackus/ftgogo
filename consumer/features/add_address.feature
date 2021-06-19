@command @address
Feature: Add Consumer Address

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Can add delivery addresses to consumers
    When I add an address with:
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
    Then I expect the command to succeed

  Scenario: Can multiple delivery addresses to consumers
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
    When I add another address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "work",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    Then I expect the command to succeed


  Scenario: Adding an address to consumers that do not exist returns an error
    When I add an address with:
    """
    {
      "ConsumerID": "b456",
      "AddressID": "home",
      "Address": {
        "Street1": "123 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "12345"
      }
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    consumer not found
    """

  Scenario: Adding an address with a duplicate id returns an error
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
    When I add another address with:
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
    Then I expect the command to fail
    And the returned error message is:
    """
    address with that identifier already exists
    """
