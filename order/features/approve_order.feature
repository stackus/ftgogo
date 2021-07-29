@command @order @approve
Feature: Order Approval

  Background: Setup Resources
    Given I have initialized the restaurant "Best Foods"
    And I have submitted an order to "Best Foods" from "Able Anders"

  Scenario: Pending orders may be approved
    Given the order to "Best Foods" from "Able Anders" is "ApprovalPending"
    When I approve the order to "Best Foods" from "Able Anders" with ticket "T123"
    Then I expect the command to succeed
    And the order to "Best Foods" from "Able Anders" is "Approved"

  Scenario: Rejected orders cannot be approved
    Given I have rejected the order to "Best Foods" from "Able Anders"
    When I approve the order to "Best Foods" from "Able Anders" with ticket "T123"
    Then I expect the command to fail
    And the returned error message is "order state does not allow action"
