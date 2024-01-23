import requests
import argparse
import time
from concurrent.futures import ThreadPoolExecutor


def send_request(url, data):
    try:
        start = time.time()
        headers = {'Content-Type': 'application/json'}
        response = requests.post(
            url=url,
            json=data,
            headers=headers,
        )
        end = time.time()
        if response.status_code == 201:
            return (0, end - start)
        else:
            return (1, end - start)
    except requests.exceptions.RequestException as e:
        return (-1, None)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-n", "--number", help="Number of requests to be sent")
    args = parser.parse_args()

    number_of_requests = int(args.number) if args.number else 1000

    url = "http://localhost:9000/api/v1/urls"
    data = {
        "expireAt": "2025-02-08T09:20:41Z",
        "url": "https://google.com"
    }

    start_time = time.time()

    with ThreadPoolExecutor(max_workers=10) as executor:
        responses = list(executor.map(
            lambda _: send_request(url, data), range(number_of_requests))
        )

    success_count = sum(1 for response in responses if response[0] == 0)
    failure_count = sum(1 for response in responses if response[0] == 1)

    total_success_time = sum(response[1] for response in responses if response[1] is not None and response[0] == 0)
    total_failure_time = sum(response[1] for response in responses if response[1] is not None and response[0] == 1)
    
    average_success_time = total_success_time / success_count if success_count > 0 else 0
    average_failure_time = total_failure_time / failure_count if failure_count > 0 else 0

    print(f"Total requests: {args.number}")
    print(f"Successful requests: {success_count}")
    print(f"Failed requests: {failure_count}")
    print(f"Average success time: {average_success_time} seconds")
    print(f"Average failure time: {average_failure_time} seconds")

if __name__ == "__main__":
    main()
