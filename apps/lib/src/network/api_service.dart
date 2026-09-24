import 'dart:convert';
import 'package:http/http.dart' as http;
import '../data/models/sms_payload.dart';

class ApiService {
  final String baseUrl;

  ApiService({this.baseUrl = "http://10.0.2.2:3000"});

  Future<bool> syncSms(String userId, String rawSms) async {
    final url = Uri.parse('$baseUrl/api/v1/sms/sync');
    final payload = SmsPayload(userId: userId, rawSms: rawSms);

    try {
      final response = await http.post(
        url,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(payload.toJson()),
      );

      return response.statusCode == 201;
    } catch (e) {
      print('Network error syncing SMS: $e');
      return false;
    }
  }
}
