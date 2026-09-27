import 'package:flutter/material.dart';
import 'package:telephony/telephony.dart';
import 'src/network/api_service.dart';
import 'src/services/sms_receiver_service.dart';

void main() {
  runApp(const FineticApp());
}

class FineticApp extends StatelessWidget {
  const FineticApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Finetic',
      theme: ThemeData(
        primarySwatch: Colors.blue,
        useMaterial3: true,
      ),
      home: const DashboardScreen(),
    );
  }
}

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {

  bool _isListening = false;
  late SmsReceiverService _smsService;

  final Telephony _telephony = Telephony.instance;

  @override
  void initState() {
    super.initState();
    final apiService = ApiService(baseUrl: 'http://10.0.2.2:8080');
    _smsService = SmsReceiverService(
      apiService: apiService,
      currentUserId: 'demo-user-id-1234',
    );
  }

  Future<void> _toggleIngestion() async {
    if (_isListening) {
      setState(() {
        _isListening = false;
      });
      debugPrint('Ingestion paused.');
      return;
    }

    final permissionsGranted = await _telephony.requestPhoneAndSmsPermissions;
    if (permissionsGranted != true) {
      debugPrint('User denied SMS permissions. Cannot ingest data.');
      return;
    }
    if (!mounted) return;

    setState(() {
      _isListening = true;
    });

    debugPrint('Started listening for M-Pesa messages...');
    _telephony.listenIncomingSms(
      onNewMessage: (SmsMessage message) {
        if (_isListening) {
          _smsService.onMessageReceived(
            message.address ?? 'UNKNOWN',
            message.body ?? '',
          );
        }
      },
      listenInBackground: false,
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Finetic Ingestion Engine'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              _isListening ? Icons.sync_rounded : Icons.sync_disabled_rounded,
              size: 80,
              color: _isListening ? Colors.green : Colors.grey,
            ),
            const SizedBox(height: 24),
            Text(
              _isListening ? 'Monitoring SMS Stream' : 'Ingestion Paused',
              style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 40),
            ElevatedButton.icon(
              onPressed: _toggleIngestion,
              icon: Icon(_isListening ? Icons.stop : Icons.play_arrow),
              label: Text(_isListening ? 'Stop Ingestion' : 'Start Ingestion'),
            ),
          ],
        ),
      ),
    );
  }
}