import json
import argparse


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('file_name', type=str, help='Имя файла для чтения')

    args = parser.parse_args()
    try:
        with open(args.file_name, 'r') as file:
            content = json.load(file)
    except FileNotFoundError:
        print(f"Файл '{args.file_name}' не найден.")
    except IOError as e:
        print(f"Ошибка при чтении файла: {e}")

    output_dict = {}
    for key, value in content.items():
        main_key, sub_key = key.split('.')
        if main_key not in output_dict:
            output_dict[main_key] = {}
        output_dict[main_key][sub_key] = value

    json_data = json.dumps(output_dict, indent=4, ensure_ascii=False)
    with open('output.json', 'w', encoding='utf-8') as json_file:
        json_file.write(json_data)

    print("Результирующий словарь успешно сохранен в файл 'output.json'.")


if __name__ == "__main__":
    main()