import pandas as pd
import os

# ファイルパス
input_file = 'H29_Habatanforstudents.xls'
output_file = 'habatan.csv'

def convert_xls_to_csv():
    if not os.path.exists(input_file):
        print(f"Error: {input_file} not found.")
        return

    try:
        # Load the excel file without header to find the correct header row
        print("Reading Excel file...")
        # xlrd is required for .xls files
        df_raw = pd.read_excel(input_file, header=None, engine='xlrd')
        
        # Find the row index that contains "番号"
        header_row = None
        for idx, row in df_raw.iterrows():
            # Check if any cell in the row contains "番号"
            if row.astype(str).str.contains("番号").any():
                header_row = idx
                break
        
        if header_row is None:
            print("Could not find header row containing '番号'. Using default header.")
            df = pd.read_excel(input_file, engine='xlrd')
        else:
            print(f"Found header at row {header_row}.")
            # Read again with the correct header
            df = pd.read_excel(input_file, header=header_row, engine='xlrd')

        # Basic cleaning: Drop rows where '番号' is NaN (if '番号' column exists)
        if '番号' in df.columns:
            df = df.dropna(subset=['番号'])
            # Convert '番号' to integer if it's numeric
            try:
                df['番号'] = df['番号'].astype(int)
            except:
                pass # Keep as is if not convertible

        # Save to CSV
        df.to_csv(output_file, index=False, encoding='utf-8-sig')
        print(f"Successfully converted {input_file} to {output_file}")
        print("First 5 rows of the converted data:")
        print(df.head())

    except Exception as e:
        print(f"An error occurred: {e}")

if __name__ == "__main__":
    convert_xls_to_csv()
