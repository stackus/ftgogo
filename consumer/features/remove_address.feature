@command @address
Feature: Remove Consumer Address

  Background: Setup a consumer
    Given I setup a consumer with:
    """
    {
      "Name": "TestConsumer"
    }
    """

  Scenario: Can remove consumer addresses
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
    When I remove an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "home",
      "Address": {
        "Street1": "456 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "67890"
      }
    }
    """
    Then I expect the command to succeed

  Scenario: Updating addresses on consumers that do not exist returns an error
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
    When I remove an address with:
    """
    {
      "ConsumerID": "b456",
      "AddressID": "home",
      "Address": {
        "Street1": "456 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "67890"
      }
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    consumer not found
    """

  Scenario: Updating addresses that do not exist returns an error
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
    When I remove an address with:
    """
    {
      "ConsumerID": "<ConsumerID>",
      "AddressID": "other",
      "Address": {
        "Street1": "456 Address St.",
        "City": "HomeTown",
        "State": "HomeState",
        "Zip": "67890"
      }
    }
    """
    Then I expect the command to fail
    And the returned error message is:
    """
    address with that identifier does not exist
    """
