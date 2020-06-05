Feature: Receive a CloudEvent

    Scenario Outline: The <consumer> receives the event.

    Given v1.0 CloudEvents Attributes:
        | key               | value                           |
        |                id | 1234-1234-1234                  |
        |              type | com.example.someevent           |
        |            source | /mycontext/subcontext           |
        |              time | 2018-04-05T03:56:24Z            |
        And JSON Data:
        """
        {"message": "Hello World!"}
        """

    When the consumer is ready
        And the event is sent to "<consumer>"

    Then the consumer got the event

    Examples:
        | consumer                                      |
        | http://localhost:8080                         |
        | http://sockeye.default.104.197.182.61.xip.io  |
        | http://observer.default.104.197.182.61.xip.io |