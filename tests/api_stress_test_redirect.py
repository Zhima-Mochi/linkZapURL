import requests
import argparse
import time
from concurrent.futures import ThreadPoolExecutor

def send_request(url):
    try:
        start = time.time()
        response = requests.get(url=url, allow_redirects=False)
        end = time.time()

        if response.status_code == 301:
            return (0, end - start)
        else:
            return (1, end - start)
    except requests.exceptions.RequestException as e:
        return (-1, None)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-c", "--code", help="Shortened code", default="abc2345")
    parser.add_argument("-n", "--number", help="Number of requests to be sent", type=int, default=1000)
    args = parser.parse_args()

    url = f"http://localhost/{args.code}"

    with ThreadPoolExecutor(max_workers=10) as executor:
        responses = list(executor.map(lambda _: send_request(url), range(args.number)))

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
